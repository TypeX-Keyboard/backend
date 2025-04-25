package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) FindBotByAddress(ctx context.Context, req *v1.FindBotByAddressReq) (res *v1.FindBotByAddressRes, err error) {
	res = &v1.FindBotByAddressRes{}
	bot, err := service.Bot().FindByAddress(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res = (*v1.FindBotByAddressRes)(service.Bot().Entity2BotInfo(ctx, bot))
	return res, gerror.NewCode(gcode.CodeOK)
}
