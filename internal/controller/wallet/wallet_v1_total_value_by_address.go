package wallet

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/wallet/v1"
)

func (c *ControllerV1) TotalValueByAddress(ctx context.Context, req *v1.TotalValueByAddressReq) (res *v1.TotalValueByAddressRes, err error) {
	res = &v1.TotalValueByAddressRes{}
	balances, err := okx.Okx().TotalValueByAddress(ctx, req.Chains, req.Address, req.AssetType, req.ExcludeRiskToken)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	if len(balances.Data) > 0 {
		*res = balances.Data
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
