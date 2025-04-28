package token

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"keyboard-api-go/internal/service"

	"keyboard-api-go/api/token/v1"
)

func (c *ControllerV1) CheckAddress(ctx context.Context, req *v1.CheckAddressReq) (res *v1.CheckAddressRes, err error) {
	res = &v1.CheckAddressRes{}
	address, err := service.Token().CheckAddress(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	*res = v1.CheckAddressRes(*address)
	return res, gerror.NewCode(gcode.CodeOK)
}
