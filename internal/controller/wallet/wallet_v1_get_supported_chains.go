package wallet

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/util/chain/okx"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/wallet/v1"
)

func (c *ControllerV1) GetSupportedChains(ctx context.Context, req *v1.GetSupportedChainsReq) (res *v1.GetSupportedChainsRes, err error) {
	res = &v1.GetSupportedChainsRes{}
	chains, err := okx.Okx().GetSupportedChains(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	*res = chains
	return res, gerror.NewCode(gcode.CodeOK)
}
