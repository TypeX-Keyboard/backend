package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolGetTxInfo(ctx context.Context, req *v1.SolGetTxInfoReq) (res *v1.SolGetTxInfoRes, err error) {
	res = &v1.SolGetTxInfoRes{}
	out, err := service.Chain().SolGetTxInfo(ctx, req.Hex)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.GetTransactionResult = out
	return res, gerror.NewCode(gcode.CodeOK)
}
