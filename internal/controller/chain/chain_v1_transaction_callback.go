package chain

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "keyboard-api-go/api/chain/v1"
	"keyboard-api-go/internal/service"
)

func (c *ControllerV1) TransactionCallback(ctx context.Context, req *v1.TransactionCallbackReq) (res *v1.TransactionCallbackRes, err error) {
	res = &v1.TransactionCallbackRes{}
	if len(req.Hash) == 0 || len(req.SignAddress) == 0 {
		return res, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	err = service.Chain().TxCallback(ctx, req.Hash, req.SignAddress)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
