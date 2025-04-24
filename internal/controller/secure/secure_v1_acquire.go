package secure

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "keyboard-api-go/api/secure/v1"
	"keyboard-api-go/internal/service"
)

func (c *ControllerV1) Acquire(ctx context.Context, req *v1.AcquireReq) (res *v1.AcquireRes, err error) {
	res = &v1.AcquireRes{}
	key, err := service.Secure().AcquirePublicKey(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.PublicKey = key
	return res, gerror.NewCode(gcode.CodeOK)
}
