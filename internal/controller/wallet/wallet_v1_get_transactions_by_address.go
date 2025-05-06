package wallet

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/okx"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "keyboard-api-go/api/wallet/v1"
)

func (c *ControllerV1) GetTransactionsByAddress(ctx context.Context, req *v1.GetTransactionsByAddressReq) (res *v1.GetTransactionsByAddressRes, err error) {
	res = &v1.GetTransactionsByAddressRes{}
	transactions, err := okx.Okx().TransactionsByAddress(ctx, req.Chains, req.Address, req.TokenAddress, req.Begin, req.End, req.Cursor, req.Limit)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	txListMap := make(map[string]model.TransactionsByAddress)
	tokensMap := make(map[string]struct{})
	tokensStrs := make([]string, 0)
	//1=接收，2=发送，3=买入，4=卖出
	for _, datum := range transactions.Data {
		res.Cursor = datum.Cursor
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
		txListMap[h] = service.Wallet().BuildHashHistory(ctx, address, req.Address)
	}
	info, err := dex.NewJup().GetPoolInfo(ctx, strings.Join(tokensStrs, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
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
	if len(txList) > 2 {
		lastTx := txList[len(txList)-1]
		res.Cursor = gconv.String(gconv.Int64(lastTx.TxTime) / 1000)
		txList = txList[:len(txList)-1]
	}
	res.TransactionList = txList
	return res, gerror.NewCode(gcode.CodeOK)
}
