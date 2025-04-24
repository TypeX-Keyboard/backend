package secure

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "keyboard-api-go/api/secure/v1"
	"keyboard-api-go/internal/service"
)

func (c *ControllerV1) Submit(ctx context.Context, req *v1.SubmitReq) (res *v1.SubmitRes, err error) {
	res = &v1.SubmitRes{}
	val := ctx.Value("uuid")
	if val == nil {
		return res, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	uuid := val.(string)
	if err := service.Secure().SubmitKey(ctx, uuid, req.AesKey); err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
