package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) SetAutoAcceptFriend(ctx context.Context, req *v1.SetAutoAcceptFriendReq) (res *v1.SetAutoAcceptFriendRes, err error) {
	res = &v1.SetAutoAcceptFriendRes{}
	err = service.Bot().UpdateAutoAccept(ctx, req.Address, req.Auto)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
