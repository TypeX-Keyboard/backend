package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) HandleFriendMsg(ctx context.Context, req *v1.HandleFriendMsgReq) (res *v1.HandleFriendMsgRes, err error) {
	res = &v1.HandleFriendMsgRes{}
	err = service.Bot().HandleFriendMsg(ctx, req.Id, req.Address, req.Accept)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
