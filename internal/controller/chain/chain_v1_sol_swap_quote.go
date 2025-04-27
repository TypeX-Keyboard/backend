package chain

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolSwapQuote(ctx context.Context, req *v1.SolSwapQuoteReq) (res *v1.SolSwapQuoteRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
