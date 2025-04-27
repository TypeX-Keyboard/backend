package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolSwapRoute(ctx context.Context, req *v1.SolSwapRouteReq) (res *v1.SolSwapRouteRes, err error) {
	res = &v1.SolSwapRouteRes{}
	swapRouter, err := service.Chain().SolSwapRoute(ctx, req.Adaptor, req.SignAddress, req.TokenIn, req.TokenOut, req.AmountIn, req.SlippageBps, req.PriorityFee, req.Mev)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.SwapRouterRes = *swapRouter
	return res, gerror.NewCode(gcode.CodeOK)
}
