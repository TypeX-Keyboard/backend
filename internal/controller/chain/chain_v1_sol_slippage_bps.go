package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolSlippageBps(ctx context.Context, req *v1.SolSlippageBpsReq) (res *v1.SolSlippageBpsRes, err error) {
	res = &v1.SolSlippageBpsRes{}
	fee, err := service.Chain().SolSlippageBps(ctx, req.TokenIn, req.TokenOut, req.AmountIn)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.ComputedAutoSlippage = fee
	return res, gerror.NewCode(gcode.CodeOK)
}
