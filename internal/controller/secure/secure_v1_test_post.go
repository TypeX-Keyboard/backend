package secure

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/secure/v1"
)

func (c *ControllerV1) TestPost(ctx context.Context, req *v1.TestPostReq) (res *v1.TestPostRes, err error) {
	res = &v1.TestPostRes{Keyword: req.Keyword}
	return res, gerror.NewCode(gcode.CodeOK)
}
