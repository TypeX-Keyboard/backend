package bot

import (
	"context"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) CreateBot(ctx context.Context, req *v1.CreateBotReq) (res *v1.CreateBotRes, err error) {
	res = &v1.CreateBotRes{}
	val := ctx.Value(consts.DeviceIdParamKey)
	if val == nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	deviceId := val.(string)
	bot, err := service.Bot().CreateBot(ctx, req.Address, deviceId)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res = (*v1.CreateBotRes)(service.Bot().Entity2BotInfo(ctx, bot))
	return res, gerror.NewCode(gcode.CodeOK)
}
