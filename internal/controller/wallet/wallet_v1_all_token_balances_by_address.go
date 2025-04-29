package wallet

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/wallet/v1"
)

func (c *ControllerV1) AllTokenBalancesByAddress(ctx context.Context, req *v1.AllTokenBalancesByAddressReq) (res *v1.AllTokenBalancesByAddressRes, err error) {
	res = &v1.AllTokenBalancesByAddressRes{}
	balances, err := okx.Okx().AllTokenBalancesByAddress(ctx, req.Chains, req.Address, req.Filter)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.Data = make([]model.TokenAsset, 0)
	if len(balances.Data) == 0 {
		return res, gerror.NewCode(gcode.CodeOK)
	}
	res.Data = service.Wallet().BuildTokenAsset(ctx, req.Address, balances.Data[0].TokenAssets)
	earnings, rate, dailyEarnings, dailyEarningsRate, err := service.Wallet().Earnings(ctx, req.Address, balances)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.GasFee = 0.000005
	res.Earnings = earnings
	res.EarningsRate = rate
	res.DailyEarnings = dailyEarnings
	res.DailyEarningsRate = dailyEarningsRate
	return res, gerror.NewCode(gcode.CodeOK)
}
