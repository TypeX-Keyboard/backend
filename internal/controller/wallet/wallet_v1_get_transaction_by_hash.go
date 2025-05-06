package wallet

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "keyboard-api-go/api/wallet/v1"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/chain/okx"
)

func (c *ControllerV1) GetTransactionByHash(ctx context.Context, req *v1.GetTransactionByHashReq) (res *v1.GetTransactionByHashRes, err error) {
	res = &v1.GetTransactionByHashRes{}
	if len(req.Address) == 0 {
		return res, gerror.NewCode(gcode.CodeInternalError, "address is empty")
	}
	detail, err := okx.Okx().TransactionDetailByTxHash(ctx, req.Chain, req.TxHash)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
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
			if transaction.From == req.Address {
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
			if transaction.To == req.Address {
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
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
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
	return res, gerror.NewCode(gcode.CodeOK)
}
