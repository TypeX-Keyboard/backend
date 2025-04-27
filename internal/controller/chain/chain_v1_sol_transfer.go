package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolTransfer(ctx context.Context, req *v1.SolTransferReq) (res *v1.SolTransferRes, err error) {
	res = &v1.SolTransferRes{}
	g.Log().Warning(ctx, req)
	transferRes, err := service.Chain().SolTransfer(ctx, req.From, req.To, req.Token, req.Amount)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.TransferRes = *transferRes
	return res, gerror.NewCode(gcode.CodeOK)
}
