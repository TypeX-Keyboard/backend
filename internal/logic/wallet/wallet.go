package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain/okx"
	"keyboard-api-go/internal/util/str"
	"sort"
	"strings"
	"time"
)

type sWallet struct {
}

func init() {
	service.RegisterWallet(New())
}

func New() service.IWallet {
	return &sWallet{}
}

// Record daily to realize benefits
func (s *sWallet) Realized(ctx context.Context) {
	_, err := gcron.Add(ctx, "@daily", func(ctx context.Context) {
		// Realized gains are recorded on a daily basis
		_, err := dao.Earning.Ctx(ctx).Where(1).Data(g.Map{
			dao.Earning.Columns().Usd24H: gdb.Raw(dao.Earning.Columns().Usd),
		}).Update()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		// Cache the price of the hold token at 0 pips
		tokens := make([]string, 0)
		type holdToken struct {
			TokenAddress string `json:"token_address"`
		}
		holds := make([]holdToken, 0)
		err = dao.Hold.Ctx(ctx).Fields(dao.Hold.Columns().TokenAddress).Distinct().Scan(&holds)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		for _, hold := range holds {
			tokens = append(tokens, hold.TokenAddress)
		}
		slice := str.SplitSlice(tokens, 100)
		for _, token := range slice {
			input := make([]okx.GetPriceInput, 0)
			for _, s := range token {
				input = append(input, okx.GetPriceInput{
					ChainIndex: 501,
					Address:    s,
				})
			}
			price, err := okx.Okx().GetPrice(ctx, input)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}
			for _, p := range price.Data {
				err := cache.GetCache().Set(ctx, fmt.Sprintf(consts.PriceIn0, p.TokenAddress), gconv.Float64(p.Price), 3*24*time.Hour)
				if err != nil {
					g.Log().Error(ctx, err)
					continue
				}
			}
		}
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
}

func (s *sWallet) SyncHold(ctx context.Context, address string, byAddress []okx.TokenAsset) (map[string]model.TokenObj, error) {
	if byAddress == nil {
		return nil, nil
	}
	holds := make([]entity.Hold, 0)
	err := dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Address, address).Scan(&holds)
	if err != nil {
		return nil, err
	}
	prices := make(map[string]model.TokenObj)
	for _, token := range byAddress {
		if token.Symbol == "SOL" && token.TokenAddress == "" {
			token.TokenAddress = consts.SolAddress
		}
		prices[token.TokenAddress] = model.TokenObj{
			Amount: gconv.Float64(token.Balance),
			Price:  gconv.Float64(token.TokenPrice),
		}
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, hold := range holds {
			token, ok := prices[hold.TokenAddress]
			if ok {
				if _, err := dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Id, hold.Id).Data(g.Map{
					dao.Hold.Columns().Amount: token.Amount,
				}).Update(); err != nil {
					g.Log().Error(ctx, err)
				}
			}
			delete(prices, hold.TokenAddress)
		}
		data := make([]entity.Hold, 0)
		for tokenAddress, token := range prices {
			price := token.Price
			_, err := cache.GetCache().SetIfNotExist(ctx, fmt.Sprintf(consts.PriceIn0, tokenAddress), price, 3*24*time.Hour)
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}
			dailyPrice, err := cache.GetCache().Get(ctx, fmt.Sprintf(consts.PriceIn0, tokenAddress))
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}
			if dailyPrice != nil && dailyPrice.Float64() != 0 {
				price = dailyPrice.Float64()
			}

			data = append(data, entity.Hold{
				Address:      address, // Wallet address
				TokenAddress: tokenAddress,
				Amount:       token.Amount, // balance
				CostPrice:    price,        // price
			})
		}
		if len(data) > 0 {
			_, err := dao.Hold.Ctx(ctx).Data(data).Insert()
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
		}
		return nil
	})
	return prices, err
}

func (s *sWallet) BuildTokenAsset(ctx context.Context, address string, balances []okx.TokenAsset) []model.TokenAsset {
	res := make([]model.TokenAsset, 0)
	_, err := s.SyncHold(ctx, address, balances)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	tokens := make(map[string]model.TokenAsset)
	tokenAddress := make([]string, 0)
	for _, asset := range balances {
		token := model.TokenAsset{
			TokenAsset: asset,
			CostPrice:  "",
			Earning:    0,
		}
		if strings.ToLower(asset.Symbol) == "sol" && token.TokenAddress == "" {
			token.TokenAddress = consts.SolAddress
		}
		tokens[token.TokenAddress] = token
		tokenAddress = append(tokenAddress, token.TokenAddress)
	}
	holdList := make([]entity.Hold, 0)
	err = dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Address, address).WhereIn(dao.Hold.Columns().TokenAddress, tokenAddress).Scan(&holdList)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	txRecords := make([]entity.TransactionRecord, 0)
	now := time.Now()
	todayStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)
	err = dao.TransactionRecord.Ctx(ctx).Where(dao.TransactionRecord.Columns().SignAddress, address).
		WhereGTE(dao.TransactionRecord.Columns().CreatedAt, todayStart).
		Where(dao.TransactionRecord.Columns().Status, 1).
		Where(dao.TransactionRecord.Columns().IsHandle, 1).Scan(&txRecords)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	type token struct {
		Amount   float64 `json:"amount"`
		UsdValue float64 `json:"usdValue"`
	}
	tokenOutMap := make(map[string]token)
	tokenInMap := make(map[string]token)
	for _, record := range txRecords {
		var usdValue float64 = 0
		if record.TokenOut == consts.SolAddress {
			usdValue = record.AmountOut * record.SolPrice
		}
		if record.TokenIn == consts.SolAddress {
			usdValue = record.AmountIn * record.SolPrice
		}
		if t, ok := tokenOutMap[record.TokenOut]; ok {
			t.Amount += record.AmountOut
			t.UsdValue += usdValue
			tokenOutMap[record.TokenOut] = t
		} else {
			tokenOutMap[record.TokenOut] = token{
				Amount:   record.AmountOut,
				UsdValue: usdValue,
			}
		}
		if t, ok := tokenInMap[record.TokenIn]; ok {
			t.Amount += record.AmountIn
			t.UsdValue += usdValue
			tokenInMap[record.TokenIn] = t
		} else {
			tokenInMap[record.TokenIn] = token{
				Amount:   record.AmountIn,
				UsdValue: usdValue,
			}
		}
	}
	for _, hold := range holdList {
		t, ok := tokens[hold.TokenAddress]
		if ok {
			t.CostPrice = fmt.Sprintf("%f", hold.CostPrice)
			price := gconv.Float64(t.TokenPrice)
			t.Earning = (price-hold.CostPrice)*hold.Amount + hold.Earning - hold.Cost
			t.EarningRate = int((t.Earning) * 10000 / (hold.CostPrice*hold.Amount + hold.Cost)) // The smallest unit is 0.01% = 1
			dailyPrice, err := cache.GetCache().Get(ctx, fmt.Sprintf(consts.PriceIn0, hold.TokenAddress))
			if err != nil {
				g.Log().Error(ctx, err)
				continue
			}
			var amount = hold.Amount
			if o, ok := tokenOutMap[hold.TokenAddress]; ok {
				if i, ok := tokenInMap[hold.TokenAddress]; ok {
					o.Amount -= i.Amount
					if o.Amount < 0 {
						o.Amount = 0
					}
					o.UsdValue -= i.UsdValue
					if o.UsdValue < 0 {
						o.UsdValue = 0
					}
					t.DailyEarning += o.UsdValue
				}
				amount = hold.Amount - o.Amount
			}
			t.DailyEarning += (dailyPrice.Float64() - hold.CostPrice) * amount
			t.DailyEarningRate = int((dailyPrice.Float64() - hold.CostPrice) * 10000 / hold.CostPrice) // The smallest unit is 0.01% = 1
			tokens[hold.TokenAddress] = t
		}
	}
	for _, token := range tokens {
		res = append(res, token)
	}
	sort.Slice(res, func(i, j int) bool {
		iPrice := gconv.Float64(res[i].TokenPrice)
		iBalance := gconv.Float64(res[i].Balance)
		jPrice := gconv.Float64(res[j].TokenPrice)
		jBalance := gconv.Float64(res[j].Balance)
		return iPrice*iBalance > jPrice*jBalance
	})
	return res
}

func (s *sWallet) Earnings(ctx context.Context, address string, byAddress *okx.AllTokenBalancesByAddressRes) (earnings float64, earningsRate int, DailyEarnings float64, dailyEarningsRate int, err error) {
	tokenAssets := make([]okx.TokenAsset, 0)
	tokens := make(map[string]model.TokenObj)
	if byAddress != nil && len(byAddress.Data) > 0 {
		tokenAssets = byAddress.Data[0].TokenAssets
		for _, token := range tokenAssets {
			if token.Symbol == "SOL" && token.TokenAddress == "" {
				token.TokenAddress = consts.SolAddress
			}
			tokens[token.TokenAddress] = model.TokenObj{
				Amount: gconv.Float64(token.Balance),
				Price:  gconv.Float64(token.TokenPrice),
			}
		}
	}
	holdList := make([]entity.Hold, 0)
	err = dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Address, address).Scan(&holdList)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// Calculate the cost of holding a position
	var holdCost float64 = 0
	// Calculate the cost for the day
	var dailyCost float64 = 0
	txRecords := make([]entity.TransactionRecord, 0)
	now := time.Now()
	todayStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)
	err = dao.TransactionRecord.Ctx(ctx).Where(dao.TransactionRecord.Columns().SignAddress, address).
		WhereGTE(dao.TransactionRecord.Columns().CreatedAt, todayStart).
		Where(dao.TransactionRecord.Columns().Status, 1).
		Where(dao.TransactionRecord.Columns().IsHandle, 1).Scan(&txRecords)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	type token struct {
		Amount   float64 `json:"amount"`
		UsdValue float64 `json:"usdValue"`
	}
	tokenOutMap := make(map[string]token)
	tokenInMap := make(map[string]token)
	for _, record := range txRecords {
		var usdValue float64 = 0
		if record.TokenOut == consts.SolAddress {
			usdValue = record.AmountOut * record.SolPrice
		}
		if record.TokenIn == consts.SolAddress {
			usdValue = record.AmountIn * record.SolPrice
		}
		if t, ok := tokenOutMap[record.TokenOut]; ok {
			t.Amount += record.AmountOut
			t.UsdValue += usdValue
			tokenOutMap[record.TokenOut] = t
		} else {
			tokenOutMap[record.TokenOut] = token{
				Amount:   record.AmountOut,
				UsdValue: usdValue,
			}
		}
		if t, ok := tokenInMap[record.TokenIn]; ok {
			t.Amount += record.AmountIn
			t.UsdValue += usdValue
			tokenInMap[record.TokenIn] = t
		} else {
			tokenInMap[record.TokenIn] = token{
				Amount:   record.AmountIn,
				UsdValue: usdValue,
			}
		}
	}
	for _, hold := range holdList {
		holdCost += hold.CostPrice * hold.Amount
		dailyPrice, err := cache.GetCache().Get(ctx, fmt.Sprintf(consts.PriceIn0, hold.TokenAddress))
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
		var amount = hold.Amount
		if o, ok := tokenOutMap[hold.TokenAddress]; ok {
			if i, ok := tokenInMap[hold.TokenAddress]; ok {
				o.Amount -= i.Amount
				if o.Amount < 0 {
					o.Amount = 0
				}
				o.UsdValue -= i.UsdValue
				if o.UsdValue < 0 {
					o.UsdValue = 0
				}
				dailyCost += o.UsdValue
			}
			amount = hold.Amount - o.Amount
		}
		dailyCost += dailyPrice.Float64() * amount
	}
	g.Log().Infof(ctx, "A list of assets %+v", tokens)
	g.Log().Infof(ctx, "A list of positions %+v", holdList)
	g.Log().Infof(ctx, "Calculate the cost of holding a position %f; Calculate the cost for the day %f", holdCost, dailyCost)
	// Calculate total assets
	var totalValue float64
	for _, datum := range tokens {
		totalValue += datum.Amount * datum.Price
	}
	// Calculate the earnings
	earnings = totalValue - holdCost
	// Calculate the yield
	earningsRate = int((earnings) * 10000 / holdCost) // The smallest unit is 0.01% = 1
	// Calculate the daily realized gain
	var earning entity.Earning
	if err := dao.Earning.Ctx(ctx).Where(dao.Earning.Columns().Address, address).Scan(&earning); err != nil && !errors.Is(err, sql.ErrNoRows) {
		g.Log().Error(ctx, err)
		return earnings, earningsRate, DailyEarnings, dailyEarningsRate, err
	}
	realizeDailyEarnings := earning.Usd - earning.Usd24H
	// Calculate daily unrealized gains
	unrealizedDailyEarnings := totalValue - dailyCost
	// Calculate daily earnings
	DailyEarnings = realizeDailyEarnings + unrealizedDailyEarnings
	// Calculate the daily yield
	dailyEarningsRate = int((DailyEarnings) * 10000 / holdCost) // 最小单位0.01%=1
	g.Log().Infof(ctx, "Calculate total assets %f; Calculate the return %f; Calculate the yield %d", totalValue, earnings, earningsRate)
	g.Log().Infof(ctx, "Calculate daily realized gain %f; Calculate the daily unrealized gain %f", realizeDailyEarnings, unrealizedDailyEarnings)
	g.Log().Infof(ctx, "Calculate daily earnings %f; Calculate the daily yield %d", DailyEarnings, dailyEarningsRate)
	return
}

func (s *sWallet) BuildHashHistory(ctx context.Context, tx model.TransactionsByAddress, walletAddress string) model.TransactionsByAddress {
	tokenInMap := make(map[string]model.TokenData)
	tokenOutMap := make(map[string]model.TokenData)
	for _, transaction := range tx.Transactions {
		from, to := "", ""
		for _, v := range transaction.From {
			from = v.Address
		}
		for _, v := range transaction.To {
			to = v.Address
		}
		if from == to {
			continue
		}
		if from == walletAddress {
			if token, ok := tokenInMap[transaction.TokenAddress]; ok {
				token.Amount += gconv.Float64(transaction.Amount)
				tokenInMap[transaction.TokenAddress] = token
			} else {
				tokenInMap[transaction.TokenAddress] = model.TokenData{
					Symbol:       transaction.Symbol,
					TokenAddress: transaction.TokenAddress,
					Amount:       gconv.Float64(transaction.Amount),
				}
			}
		}
		if to == walletAddress {
			if token, ok := tokenOutMap[transaction.TokenAddress]; ok {
				token.Amount += gconv.Float64(transaction.Amount)
				tokenOutMap[transaction.TokenAddress] = token
			} else {
				tokenOutMap[transaction.TokenAddress] = model.TokenData{
					Symbol:       transaction.Symbol,
					TokenAddress: transaction.TokenAddress,
					Amount:       gconv.Float64(transaction.Amount),
				}
			}
		}
	}
	action := 0
	if len(tokenOutMap) > 0 {
		action = 1
	}
	if len(tokenInMap) > 0 {
		action = 2
	}
	for s, f := range tokenOutMap {
		// It is SOL that is received
		if s == consts.SolAddress && len(tokenInMap) > 0 {
			action = 4
		}
		tx.TokenAddressOut = s
		tx.OutAmount = gconv.String(f.Amount)
		tx.SymbolOut = f.Symbol
	}
	for s, f := range tokenInMap {
		// Exclude pinch-proof transfers
		if len(tokenInMap) > 1 && s == consts.SolAddress {
			continue
		}
		// It's sol that's sending
		if s == consts.SolAddress && len(tokenOutMap) > 0 {
			action = 3
		}
		tx.TokenAddressIn = s
		tx.InAmount = gconv.String(f.Amount)
		tx.SymbolIn = f.Symbol
	}
	if action <= 2 {
		if len(tokenInMap) > 0 {
			tx.TokenAddress = tx.TokenAddressIn
			tx.Amount = tx.InAmount
			tx.Symbol = tx.SymbolIn
		}
		if len(tokenOutMap) > 0 {
			tx.TokenAddress = tx.TokenAddressOut
			tx.Amount = tx.OutAmount
			tx.Symbol = tx.SymbolOut
		}
	}
	tx.Action = action
	return tx
}
