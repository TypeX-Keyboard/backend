package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolPriorityFee(ctx context.Context, req *v1.SolPriorityFeeReq) (res *v1.SolPriorityFeeRes, err error) {
	res = &v1.SolPriorityFeeRes{}
	fee, err := service.Chain().SolPriorityFee(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.PriorityFeeRes = fee
	return res, gerror.NewCode(gcode.CodeOK)
}
