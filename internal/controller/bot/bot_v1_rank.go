package bot

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/bot/v1"
)

func (c *ControllerV1) Rank(ctx context.Context, req *v1.RankReq) (res *v1.RankRes, err error) {
	res = &v1.RankRes{}
	rank, selfRank, amount, err := service.Bot().Rank(ctx, req.Address)
	if err != nil {
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.List = rank
	res.SelfRank = selfRank
	res.Amount = amount
	return res, gerror.NewCode(gcode.CodeOK)
}
