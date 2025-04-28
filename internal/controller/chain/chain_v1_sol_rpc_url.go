package chain

import (
	"context"
	"keyboard-api-go/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/chain/v1"
)

func (c *ControllerV1) SolRpcURL(ctx context.Context, req *v1.SolRpcURLReq) (res *v1.SolRpcURLRes, err error) {
	res = &v1.SolRpcURLRes{}
	res.URL = service.Chain().SolRpcURL(ctx)
	res.MevURL = "https://mainnet.block-engine.jito.wtf/api/v1/transactions"
	return res, gerror.NewCode(gcode.CodeOK)
}
