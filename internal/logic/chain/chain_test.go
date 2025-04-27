package chain

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/service"
	"testing"
	"time"
)

func init() {
	service.RegisterChain(New())
}

var ctx = gctx.New()

func TestSChain_WeightedAverage(t *testing.T) {
	average, err := service.Chain().WeightedAverage([]float64{10, 1}, []float64{2, 0})
	if err != nil {
		t.Error(err)
	}
	g.Dump(average)
	created := time.Now().Add(-32 * time.Minute)
	sub := time.Now().Sub(created)
	g.Dump(sub > 30*time.Minute)
}

func TestPushList(t *testing.T) {
	hash := "61cUCoc2Ayxtgw6X82DTzE7vTQZVxS1nMWyD1zjpXeL7BsW9oEA8HMJsxwJZeodDLcmJyQ9vuoabVv42wzypqtn2"
	_, err := g.Redis().LPush(ctx, consts.OKXGetTXList, hash)
	if err != nil {
		g.Log().Error(ctx, err)
	}
}

func TestPushSuccessList(t *testing.T) {
	hash := "3GYzDUw6j1BnGYdUQqohkVuMqVdaDHvByjrU77GXW2jZEDcTfHeLnzeGnYGM1p9g8UMhCgFzChx81ZsqtVbw3yCz"
	_, err := g.Redis().LPush(ctx, consts.TxSuccessHandle, hash)
	if err != nil {
		g.Log().Error(ctx, err)
	}
}
