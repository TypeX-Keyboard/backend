package token

import (
	"context"
	"keyboard-api-go/internal/util/chain/solana/dex"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/token/v1"
)

func (c *ControllerV1) GetTokenDetail(ctx context.Context, req *v1.GetTokenDetailReq) (res *v1.GetTokenDetailRes, err error) {
	res = &v1.GetTokenDetailRes{}
	detail, err := dex.NewJup().GetPoolInfo(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	for _, pool := range detail.Pools {
		*res = v1.GetTokenDetailRes(pool)
	}
	return res, gerror.NewCode(gcode.CodeOK)
}
