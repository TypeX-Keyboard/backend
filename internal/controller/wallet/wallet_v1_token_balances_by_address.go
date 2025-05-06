package wallet

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/wallet/v1"
)

func (c *ControllerV1) TokenBalancesByAddress(ctx context.Context, req *v1.TokenBalancesByAddressReq) (res *v1.TokenBalancesByAddressRes, err error) {
	res = &v1.TokenBalancesByAddressRes{}
	balances, err := okx.Okx().TokenBalancesByAddress(ctx, req.Chains, req.Address, req.TokenAddresses)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	if len(balances.Data) == 0 {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	*res = service.Wallet().BuildTokenAsset(ctx, req.Address, balances.Data[0].TokenAssets)
	return res, gerror.NewCode(gcode.CodeOK)
}
