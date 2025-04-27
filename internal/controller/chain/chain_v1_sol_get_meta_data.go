package chain

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolGetMetaData(ctx context.Context, req *v1.SolGetMetaDataReq) (res *v1.SolGetMetaDataRes, err error) {
	res = &v1.SolGetMetaDataRes{}
	data, err := service.Chain().SolGetMetaData(ctx, req.Adaptor, req.Tokens)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}
	res.TokenInfoRes = data
	return res, gerror.NewCode(gcode.CodeOK)
}
