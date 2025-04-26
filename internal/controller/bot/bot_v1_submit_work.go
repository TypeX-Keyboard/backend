package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) SubmitWork(ctx context.Context, req *v1.SubmitWorkReq) (res *v1.SubmitWorkRes, err error) {
	res = &v1.SubmitWorkRes{}
	err = service.Bot().SubmitWork(ctx, req.Address, req.TypeCount)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
