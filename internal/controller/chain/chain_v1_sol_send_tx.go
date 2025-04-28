package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolSendTx(ctx context.Context, req *v1.SolSendTxReq) (res *v1.SolSendTxRes, err error) {
	res = &v1.SolSendTxRes{}
	tx, err := service.Chain().SolSendTx(ctx, req.TxRaw, req.SignAddress)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.Hex = tx
	return res, gerror.NewCode(gcode.CodeOK)
}
