package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) FriendList(ctx context.Context, req *v1.FriendListReq) (res *v1.FriendListRes, err error) {
	res = &v1.FriendListRes{}
	list, err := service.Bot().FriendList(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	*res = list
	return res, gerror.NewCode(gcode.CodeOK)
}
