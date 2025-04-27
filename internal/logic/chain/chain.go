package chain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/dao"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain/okx"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"strings"
	"time"
)

type sChain struct {
}

func init() {
	service.RegisterChain(New())
}

func New() service.IChain {
	return &sChain{}
}

func (s *sChain) SyncSolPrice(ctx context.Context) {
	_, _ = gcron.Add(ctx, "*/5 * * * * *", func(ctx context.Context) {
		info, err := dex.NewJup().GetPoolInfo(ctx, consts.SolAddress)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		for _, pool := range info.Pools {
			err := cache.GetCache().Set(ctx, consts.SolPrice, pool.BaseAsset.UsdPrice, 0)
			if err != nil {
				g.Log().Error(ctx, err)
			}
		}
	})
}

// TxCallback 函数用于处理交易回调
func (s *sChain) TxCallback(ctx context.Context, hash, SignAddress string) error {
	var solPrice float64 = 0
	get, err := cache.GetCache().Get(ctx, consts.SolPrice)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if get != nil {
		solPrice = get.Float64()
	}
	_, err = dao.TransactionRecord.Ctx(ctx).Data(entity.TransactionRecord{
		Hash:        hash,
		SignAddress: SignAddress,
		SolPrice:    solPrice,
	}).Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	go func(ctx context.Context, hash string) {
		time.Sleep(5 * time.Second)
		_, err = g.Redis().LPush(ctx, consts.OKXGetTXList, hash)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		g.Log().Warningf(ctx, "%s push TxListen queue", hash)
	}(gctx.New(), hash)

	return err
}

func (s *sChain) TxListen(ctx context.Context) {
	s.SyncTx(ctx)
	ch := make(chan struct{}, 10)
	defer close(ch)
	for {
		members, err := g.Redis().BRPop(ctx, 10, consts.OKXGetTXList)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		if len(members) != 2 {
			continue
		}
		hash := members.Strings()[1]
		g.Log().Warningf(ctx, "%s pop TxListen queue", hash)
		ch <- struct{}{}
		go func(ctx context.Context, hash string) {
			defer func() {
				g.Log().Warningf(ctx, "%s TxListen queue exec finish", hash)
				<-ch
			}()
			var record *entity.TransactionRecord
			if err := dao.TransactionRecord.Ctx(ctx).
				Where(dao.TransactionRecord.Columns().Hash, hash).
				Where(dao.TransactionRecord.Columns().Status, 0).
				Scan(&record); err != nil {
				g.Log().Error(ctx, err)
				return
			}
			if record == nil {
				g.Log().Warningf(ctx, "OKXGetTXList txHash: %v record not found", hash)
				return
			}
			if record.SolPrice == 0 {
				solPrice, err := cache.GetCache().Get(ctx, consts.SolPrice)
				if err != nil {
					g.Log().Error(ctx, err)
				}
				if solPrice != nil {
					record.SolPrice = solPrice.Float64()
				}
			}
			txHash, err := okx.Okx().TransactionDetailByTxHash(ctx, "501", record.Hash)
			if err != nil {
				g.Log().Error(ctx, err)
				time.Sleep(5 * time.Second)
				// 重新入队
				_, err = g.Redis().LPush(ctx, consts.OKXGetTXList, record.Hash)
				if err != nil {
					g.Log().Error(ctx, err)
				}
				return
			}
			g.Log().Warningf(ctx, "txHash: %v", txHash)
			g.Log().Warningf(ctx, "record: %v", record)
			//if time.Now().Sub(record.CreatedAt.Time) > 30*time.Minute {
			//	record.Status = 2
			//}
			for _, tx := range txHash.Data {
				if tx.TxStatus == "success" {
					record.Status = 1
				}
				if tx.TxStatus == "fail" {
					record.Status = 2
				}
				tokenInMap := make(map[string]model.TokenData)
				tokenOutMap := make(map[string]model.TokenData)
				for _, detail := range tx.TokenTransferDetails {
					if strings.ToLower(detail.Symbol) == "sol" && detail.TokenContractAddress == "" {
						detail.TokenContractAddress = consts.SolAddress
					}
					if detail.From == detail.To {
						continue
					}
					if detail.From == record.SignAddress {
						if token, ok := tokenInMap[detail.TokenContractAddress]; ok {
							token.Amount += gconv.Float64(detail.Amount)
							tokenInMap[detail.TokenContractAddress] = token
						} else {
							tokenInMap[detail.TokenContractAddress] = model.TokenData{
								Symbol:       detail.Symbol,
								TokenAddress: detail.TokenContractAddress,
								Amount:       gconv.Float64(detail.Amount),
							}
						}
					}
					if detail.To == record.SignAddress {
						if token, ok := tokenOutMap[detail.TokenContractAddress]; ok {
							token.Amount += gconv.Float64(detail.Amount)
							tokenOutMap[detail.TokenContractAddress] = token
						} else {
							tokenOutMap[detail.TokenContractAddress] = model.TokenData{
								Symbol:       detail.Symbol,
								TokenAddress: detail.TokenContractAddress,
								Amount:       gconv.Float64(detail.Amount),
							}
						}
					}
				}
				action := 0
				if len(tokenOutMap) > 0 {
					action = 1
				}
				if len(tokenInMap) > 0 {
					action = 1
				}
				for s, f := range tokenOutMap {
					// it is sol that is received
					if s == consts.SolAddress && len(tokenInMap) > 0 {
						action = 2
					}
					record.TokenOut = s
					record.AmountOut = f.Amount
				}
				for s, f := range tokenInMap {
					// exclude pinch proof transfers
					if len(tokenInMap) > 1 && s == consts.SolAddress {
						continue
					}
					// it s sol that s sending
					if s == consts.SolAddress && len(tokenOutMap) > 0 {
						action = 2
					}
					record.TokenIn = s
					record.AmountIn = f.Amount
				}
				g.Log().Warningf(ctx, "input token: %+v", tokenInMap)
				g.Log().Warningf(ctx, "output token: %+v", tokenOutMap)
				record.Action = action
			}
			if _, err := dao.TransactionRecord.Ctx(ctx).Save(record); err != nil {
				g.Log().Error(ctx, err)
			}
			g.Log().Warningf(ctx, "record Status: %v", record.Status)
			if record.Status == 0 {
				// re enlistment
				g.Log().Warningf(ctx, "%s push TxListen queue again", hash)
				_, err = g.Redis().LPush(ctx, consts.OKXGetTXList, record.Hash)
				if err != nil {
					g.Log().Error(ctx, err)
				}
				return
			}
			_, err = g.Redis().LPush(ctx, consts.TxSuccessHandle, record.Hash)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			g.Log().Warningf(ctx, "push TxSuccessHandle record: %v", record)
			return
		}(gctx.New(), hash)
	}
}

func (s *sChain) SyncTx(ctx context.Context) {
	var list []entity.TransactionRecord
	err := dao.TransactionRecord.Ctx(ctx).Where(dao.TransactionRecord.Columns().Status, 0).Scan(&list)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	for _, record := range list {
		_, err = g.Redis().LPush(ctx, consts.OKXGetTXList, record.Hash)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
	}
}

func (s *sChain) TxSuccessListen(ctx context.Context) {
	s.SyncTxHandle(ctx)
	ch := make(chan struct{}, 10)
	defer close(ch)
	for {
		members, err := g.Redis().BRPop(ctx, 10, consts.TxSuccessHandle)
		if err != nil {
			// 处理错误
			g.Log().Error(ctx, err)
		}
		if len(members.Strings()) != 2 {
			continue
		}
		hash := members.Strings()[1]
		g.Log().Warningf(ctx, "%s pop TxSuccessHandle queue", hash)
		ch <- struct{}{}
		go func(hash string) {
			defer func() {
				g.Log().Warningf(ctx, "%s TxSuccessHandle queue exec finish", hash)
				<-ch
			}()
			var record *entity.TransactionRecord
			if err := dao.TransactionRecord.Ctx(ctx).
				Where(dao.TransactionRecord.Columns().Hash, hash).
				Where(dao.TransactionRecord.Columns().Status, 1).
				Where(dao.TransactionRecord.Columns().IsHandle, false).
				Scan(&record); err != nil {
				g.Log().Error(ctx, err)
				return
			}
			if record == nil {
				g.Log().Warningf(ctx, "TxSuccessHandle txHash: %v record not found", hash)
				return
			}
			g.Log().Warningf(ctx, "pop TxSuccessHandle record: %v", record)
			if err := s.Hold(ctx, record); err != nil {
				_, err = g.Redis().LPush(ctx, consts.TxSuccessHandle, record.Hash)
				if err != nil {
					g.Log().Error(ctx, err)
				}
				return
			}
		}(hash)
	}
}

func (s *sChain) SyncTxHandle(ctx context.Context) {
	var list []entity.TransactionRecord
	err := dao.TransactionRecord.Ctx(ctx).Where(dao.TransactionRecord.Columns().Status, 1).
		Where(dao.TransactionRecord.Columns().IsHandle, false).Scan(&list)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	for _, record := range list {
		_, err = g.Redis().LPush(ctx, consts.TxSuccessHandle, record.Hash)
		if err != nil {
			g.Log().Error(ctx, err)
			continue
		}
	}
}

func (s *sChain) Hold(ctx context.Context, txRecord *entity.TransactionRecord) error {
	input := []okx.GetPriceInput{
		{
			ChainIndex: 501,
			Address:    txRecord.TokenIn,
		},
		{
			ChainIndex: 501,
			Address:    txRecord.TokenOut,
		},
	}
	price, err := okx.Okx().GetPrice(ctx, input)
	if err != nil || price == nil {
		return err
	}
	tokens := make(map[string]float64)
	for _, token := range price.Data {
		tokens[token.TokenAddress] = gconv.Float64(token.Price)
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, ok := tokens[txRecord.TokenOut]; ok && len(txRecord.TokenOut) != 0 {
			var buy entity.Hold
			err := dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Address, txRecord.SignAddress).Where(dao.Hold.Columns().TokenAddress, txRecord.TokenOut).Scan(&buy)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				g.Log().Error(ctx, err)
				return err
			}
			// Calculate the average purchase price
			buyCost, err := s.WeightedAverage(
				[]float64{buy.CostPrice, tokens[txRecord.TokenOut]},
				[]float64{buy.Amount, txRecord.AmountOut},
			)
			if err != nil {
				g.Log().Error(ctx, err)
				return err
			}
			if buy.Id == 0 {
				buy = entity.Hold{
					Address:      txRecord.SignAddress,
					TokenAddress: txRecord.TokenOut,
					Amount:       txRecord.AmountOut,
					CostPrice:    buyCost, // Buy-in cost
				}
				if _, err := dao.Hold.Ctx(ctx).Data(buy).Insert(); err != nil {
					return err
				}
			} else {
				if _, err := dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Id, buy.Id).Data(g.Map{
					dao.Hold.Columns().Amount:    gdb.Raw(fmt.Sprintf("%s+%f", dao.Hold.Columns().Amount, txRecord.AmountOut)), // 累加
					dao.Hold.Columns().CostPrice: buyCost,                                                                      // 买入成本
				}).Update(); err != nil {
					g.Log().Error(ctx, err)
					return err
				}
			}
		}
		// Sell
		var sell entity.Hold
		err = dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Address, txRecord.SignAddress).Where(dao.Hold.Columns().TokenAddress, txRecord.TokenIn).Scan(&sell)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			g.Log().Error(ctx, err)
			return err
		}
		if sell.Id == 0 {
			sell.Amount = txRecord.AmountIn
			sell.CostPrice = tokens[txRecord.TokenIn]
		}
		total := tokens[txRecord.TokenOut] * txRecord.AmountOut                // Total value
		usd := (tokens[txRecord.TokenIn] - sell.CostPrice) * txRecord.AmountIn // Realize profit and loss
		g.Log().Warningf(ctx, "sell: %v", sell)
		g.Log().Warningf(ctx, "TokenIn price: %v,Cost price: %v,amount: %v", tokens[txRecord.TokenIn], sell.CostPrice, txRecord.AmountIn)
		g.Log().Warningf(ctx, "total: %v, usd: %v", total, usd)
		var earning entity.Earning
		err = dao.Earning.Ctx(ctx).Where(dao.Earning.Columns().Address, txRecord.SignAddress).Scan(&earning)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			g.Log().Error(ctx, err)
			return err
		}
		if sell.Id != 0 {
			if _, err := dao.Hold.Ctx(ctx).Where(dao.Hold.Columns().Id, sell.Id).Data(g.Map{
				dao.Hold.Columns().Cost:    gdb.Raw(fmt.Sprintf("%s+%f", dao.Hold.Columns().Cost, txRecord.AmountIn*sell.CostPrice)),
				dao.Hold.Columns().Earning: gdb.Raw(fmt.Sprintf("%s+%f", dao.Hold.Columns().Earning, usd)),
				dao.Hold.Columns().Amount:  gdb.Raw(fmt.Sprintf("%s-%f", dao.Hold.Columns().Amount, txRecord.AmountIn)),
			}).Update(); err != nil {
				g.Log().Error(ctx, err)
			}
		}
		// Update investment costs and realized profit and loss
		if earning.Id == 0 {
			earning = entity.Earning{
				Address: txRecord.SignAddress,
				Usd:     usd,
				Cost:    total,
			}
			if _, err := dao.Earning.Ctx(ctx).Data(earning).Insert(); err != nil {
				g.Log().Error(ctx, err)
				return err
			}
		} else {
			if _, err := dao.Earning.Ctx(ctx).Where(dao.Earning.Columns().Id, earning.Id).Data(g.Map{
				dao.Earning.Columns().Cost: gdb.Raw(fmt.Sprintf("%s+%f", dao.Earning.Columns().Cost, total)),
				dao.Earning.Columns().Usd:  gdb.Raw(fmt.Sprintf("%s+%f", dao.Earning.Columns().Usd, usd)),
			}).Update(); err != nil {
				g.Log().Error(ctx, err)
				return err
			}
		}
		_, err = dao.TransactionRecord.Ctx(ctx).Where(dao.TransactionRecord.Columns().Id, txRecord.Id).Data(g.Map{
			dao.TransactionRecord.Columns().IsHandle: true,
		}).Update()
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
		return nil
	})
	return err
}

// WeightedAverage Calculate the weighted average
func (s *sChain) WeightedAverage(values []float64, weights []float64) (float64, error) {
	if len(values) != len(weights) || len(values) == 0 {
		return 0, errors.New("The values and weights array must be equal in length and non-null")
	}

	var sumValues, sumWeights float64

	for i := 0; i < len(values); i++ {
		sumValues += values[i] * weights[i] // Value × weight
		sumWeights += weights[i]            // Cumulative weights
	}

	if sumWeights == 0 {
		return 0, errors.New("The total weight cannot be zero")
	}

	return sumValues / sumWeights, nil
}
