package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) AddFriend(ctx context.Context, req *v1.AddFriendReq) (res *v1.AddFriendRes, err error) {
	res = &v1.AddFriendRes{}
	isAutoAccept, err := service.Bot().AddFriend(ctx, req.Address, req.FriendAddress)
	if err != nil {
		return res, err
	}
	res.IsAutoAccept = isAutoAccept
	return res, gerror.NewCode(gcode.CodeOK)
}
