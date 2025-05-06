package token

import (
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/chain/solana"
	"testing"
)

var ctx = gctx.New()

func init() {
	var conf = solana.SolanaConf{}
	ctx := gctx.GetInitCtx()
	if err := solana.InitSolana(ctx, conf); err != nil {
		panic(err)
	}
}

func TestSToken_CheckAddress(t *testing.T) {
	address, err := service.Token().CheckAddress(ctx, "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN")
	if err != nil {
		t.Error(err)
	}
	t.Log(address)
}
