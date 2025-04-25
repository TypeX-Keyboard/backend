package bot

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) Active(ctx context.Context, req *v1.ActiveReq) (res *v1.ActiveRes, err error) {
	res = &v1.ActiveRes{}
	val := ctx.Value(consts.DeviceIdParamKey)
	if val == nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	deviceId := val.(string)
	_, err = service.Bot().CreateBot(ctx, req.Address, deviceId)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	err = service.Bot().Active(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
