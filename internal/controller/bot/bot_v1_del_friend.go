package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) DelFriend(ctx context.Context, req *v1.DelFriendReq) (res *v1.DelFriendRes, err error) {
	res = &v1.DelFriendRes{}
	err = service.Bot().DelFriend(ctx, req.Address, req.FriendAddress)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
