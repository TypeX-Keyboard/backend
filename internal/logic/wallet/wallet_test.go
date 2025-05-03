package wallet

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "keyboard-api-go/api/wallet/v1"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/okx"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"keyboard-api-go/internal/util/str"
	"sort"
	"strings"
	"testing"
)

func init() {
	service.RegisterWallet(New())
}

var ctx = gctx.New()

func TestRealized(t *testing.T) {
	// 每日记录已实现收益
	_, err := dao.Earning.Ctx(ctx).Where(1).Data(g.Map{
		dao.Earning.Columns().Usd24H: gdb.Raw(dao.Earning.Columns().Usd),
	}).Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	m, err := dao.Hold.Ctx(ctx).Fields(dao.Hold.Columns().TokenAddress).Distinct().All()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(m)
	type holdToken struct {
		TokenAddress string `json:"token_address"`
	}
	// 缓存hold代币0点价格
	holds := make([]holdToken, 0)
	err = dao.Hold.Ctx(ctx).Fields(dao.Hold.Columns().TokenAddress).Distinct().Scan(&holds)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(holds)
	tokens := make([]string, 0)
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
		g.Dump(input)
		//price, err := okx.Okx().GetPrice(ctx, input)
		//if err != nil {
		//	g.Log().Error(ctx, err)
		//	continue
		//}
		//for _, p := range price.Data {
		//	err := cache.GetCache().Set(ctx, fmt.Sprintf(consts.PriceIn0, p.TokenAddress), gconv.Float64(p.Price), 3*24*time.Hour)
		//	if err != nil {
		//		g.Log().Error(ctx, err)
		//		continue
		//	}
		//}
	}
}

func TestTxHistory(t *testing.T) {
	wallet := "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi"
	transactions, err := okx.Okx().TransactionsByAddress(ctx, "501", wallet, "", "", "", "", "")
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	txListMap := make(map[string]model.TransactionsByAddress)
	tokensMap := make(map[string]struct{})
	tokensStrs := make([]string, 0)
	//1=接收，2=发送，3=买入，4=卖出
	for _, datum := range transactions.Data {

		for _, transaction := range datum.TransactionList {
			if transaction.TokenAddress == "" && strings.ToLower(transaction.Symbol) == "sol" {
				transaction.TokenAddress = consts.SolAddress
			}
			if _, ok := tokensMap[transaction.TokenAddress]; !ok {
				tokensMap[transaction.TokenAddress] = struct{}{}
				tokensStrs = append(tokensStrs, transaction.TokenAddress)
			}
			if tx, ok := txListMap[transaction.TxHash]; ok {
				if transaction.TxStatus == "fail" {
					tx.TxStatus = transaction.TxStatus
				}
				tx.Transactions = append(tx.Transactions, model.TxData{
					From:         transaction.From,
					To:           transaction.To,
					TokenAddress: transaction.TokenAddress,
					Amount:       transaction.Amount,
					Symbol:       transaction.Symbol,
					TxStatus:     transaction.TxStatus,
					HitBlacklist: false,
				})
				txListMap[transaction.TxHash] = tx
			} else {
				txListMap[transaction.TxHash] = model.TransactionsByAddress{
					ChainIndex: transaction.ChainIndex,
					TxHash:     transaction.TxHash,
					TxTime:     transaction.TxTime,
					TxStatus:   transaction.TxStatus,
					Transactions: []model.TxData{
						{
							From:         transaction.From,
							To:           transaction.To,
							TokenAddress: transaction.TokenAddress,
							Amount:       transaction.Amount,
							Symbol:       transaction.Symbol,
							TxStatus:     transaction.TxStatus,
							HitBlacklist: transaction.HitBlacklist,
						},
					},
				}
			}
		}
	}
	for h, address := range txListMap {
		txListMap[h] = service.Wallet().BuildHashHistory(ctx, address, wallet)
	}
	info, err := dex.NewJup().GetPoolInfo(ctx, strings.Join(tokensStrs, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	iconMap := make(map[string]string)
	for _, pool := range info.Pools {
		iconMap[pool.BaseAsset.ID] = pool.BaseAsset.Icon
	}
	txList := make([]model.TransactionsByAddress, 0)
	for _, tx := range txListMap {
		if icon, ok := iconMap[tx.TokenAddress]; ok {
			tx.TokenLogoUrl = icon
		}
		if icon, ok := iconMap[tx.TokenAddressIn]; ok {
			tx.InTokenLogoUrl = icon
		}
		if icon, ok := iconMap[tx.TokenAddressOut]; ok {
			tx.OutTokenLogoUrl = icon
		}
		txList = append(txList, tx)
	}
	sort.Slice(txList, func(i, j int) bool {
		return txList[i].TxTime > txList[j].TxTime
	})

	for _, address := range txList {
		g.Dump(address)
	}
}

func TestTxByHash(t *testing.T) {
	res := &v1.GetTransactionByHashRes{}
	wallet := "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi"
	hash := "4UceNPdP3qi7AHjGuBsENqLP1Dc8jGGj3dUJ6HySvEudMNCto5q2XoDCRinY8tn6YpS1kAmww1MipRbVLbLoUkY2"
	detail, err := okx.Okx().TransactionDetailByTxHash(ctx, "501", hash)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	tokensMap := make(map[string]struct{})
	tokensStrs := make([]string, 0)
	//1=接收，2=发送，3=买入，4=卖出
	for _, v := range detail.Data {
		res.TxHash = v
		tokenInMap := make(map[string]model.TokenData)
		tokenOutMap := make(map[string]model.TokenData)
		for _, transaction := range v.TokenTransferDetails {
			if transaction.TokenContractAddress == "" && strings.ToLower(transaction.Symbol) == "sol" {
				transaction.TokenContractAddress = consts.SolAddress
			}
			if transaction.From == transaction.To {
				continue
			}
			// 获取平台费
			if transaction.To == consts.FeeAccount || transaction.To == consts.FeeAccountATA {
				res.PlatformFee = transaction.Amount
				continue
			}
			if transaction.From == wallet {
				if token, ok := tokenInMap[transaction.TokenContractAddress]; ok {
					token.Amount += gconv.Float64(transaction.Amount)
					tokenInMap[transaction.TokenContractAddress] = token
				} else {
					tokenInMap[transaction.TokenContractAddress] = model.TokenData{
						Symbol:       transaction.Symbol,
						TokenAddress: transaction.TokenContractAddress,
						Amount:       gconv.Float64(transaction.Amount),
					}
				}
			}
			if transaction.To == wallet {
				if token, ok := tokenOutMap[transaction.TokenContractAddress]; ok {
					token.Amount += gconv.Float64(transaction.Amount)
					tokenOutMap[transaction.TokenContractAddress] = token
				} else {
					tokenOutMap[transaction.TokenContractAddress] = model.TokenData{
						Symbol:       transaction.Symbol,
						TokenAddress: transaction.TokenContractAddress,
						Amount:       gconv.Float64(transaction.Amount),
					}
				}
			}
			if _, ok := tokensMap[transaction.TokenContractAddress]; !ok {
				tokensMap[transaction.TokenContractAddress] = struct{}{}
				tokensStrs = append(tokensStrs, transaction.TokenContractAddress)
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
			// 接收的是sol
			if s == consts.SolAddress && len(tokenInMap) > 0 {
				action = 4
			}
			res.TokenAddressOut = s
			res.OutAmount = gconv.String(f.Amount)
			res.SymbolOut = f.Symbol
		}
		for s, f := range tokenInMap {
			// 排除防夹转账
			if len(tokenInMap) > 1 && s == consts.SolAddress {
				continue
			}
			// 发送的是sol
			if s == consts.SolAddress && len(tokenOutMap) > 0 {
				action = 3
			}
			res.TokenAddressIn = s
			res.InAmount = gconv.String(f.Amount)
			res.SymbolIn = f.Symbol
		}
		if action <= 2 {
			if len(tokenInMap) > 0 {
				res.TokenAddress = res.TokenAddressIn
				res.Amount = res.InAmount
				res.Symbol = res.SymbolIn
			}
			if len(tokenOutMap) > 0 {
				res.TokenAddress = res.TokenAddressOut
				res.Amount = res.OutAmount
				res.Symbol = res.SymbolOut
			}
		}
		res.Action = action
	}
	info, err := dex.NewJup().GetPoolInfo(ctx, strings.Join(tokensStrs, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	iconMap := make(map[string]string)
	for _, pool := range info.Pools {
		iconMap[pool.BaseAsset.ID] = pool.BaseAsset.Icon
	}
	if icon, ok := iconMap[res.TokenAddress]; ok {
		res.TokenLogoUrl = icon
	}
	if icon, ok := iconMap[res.TokenAddressIn]; ok {
		res.InTokenLogoUrl = icon
	}
	if icon, ok := iconMap[res.TokenAddressOut]; ok {
		res.OutTokenLogoUrl = icon
	}
	g.Dump(res)
}
