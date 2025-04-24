package secure

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/secure/v1"
)

func (c *ControllerV1) GenKey(ctx context.Context, req *v1.GenKeyReq) (res *v1.GenKeyRes, err error) {
	res = &v1.GenKeyRes{}
	err = service.Secure().GenClientKey(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
