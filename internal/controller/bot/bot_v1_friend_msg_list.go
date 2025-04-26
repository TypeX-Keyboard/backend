package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) FriendMsgList(ctx context.Context, req *v1.FriendMsgListReq) (res *v1.FriendMsgListRes, err error) {
	res = &v1.FriendMsgListRes{}
	list, total, err := service.Bot().FriendMsgList(ctx, req.Address, req.IsSelf, req.Page, req.Size)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.List = list
	res.Total = total
	return res, gerror.NewCode(gcode.CodeOK)
}
