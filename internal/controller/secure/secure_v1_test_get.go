package secure

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/secure/v1"
)

func (c *ControllerV1) TestGet(ctx context.Context, req *v1.TestGetReq) (res *v1.TestGetRes, err error) {
	res = &v1.TestGetRes{
		Keyword: req.Keyword,
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
