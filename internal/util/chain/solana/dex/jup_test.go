package dex

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/consts"
	utilCache "keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain/solana/model"
	"testing"
)

func init() {
	var (
		cache = gcache.New()
		redis = g.Redis()
	)
	// Create redis cache adapter and set it to cache object.
	cache.SetAdapter(gcache.NewAdapterRedis(redis))
	utilCache.InitCache(gctx.GetInitCtx(), cache)
}

func TestSwapQuote(t *testing.T) {
	ctx := gctx.GetInitCtx()
	jup := NewJup()
	res, err := jup.GetSwapQuoteEarly(ctx, model.SwapQuoteReq{
		QuoteBaseReq: model.QuoteBaseReq{
			InputMint:  "So11111111111111111111111111111111111111112",
			OutputMint: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
			Amount:     100000000,
		},
		FeeAccount:  consts.FeeAccountATA,
		SlippageBps: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
	mev := false
	swapInput := model.GetSwapRawReq{
		QuoteResponse:           res,
		FeeAccount:              consts.FeeAccountATA,
		UserPublicKey:           "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi",
		UseSharedAccounts:       true,
		DynamicComputeUnitLimit: true,
		WrapAndUnwrapSol:        true,
		DynamicSlippage:         true,
	}
	if mev {
		swapInput.PrioritizationFeeLamports = map[string]interface{}{
			"jitoTipLamports": 300000,
		}
	} else {
		swapInput.PrioritizationFeeLamports = map[string]interface{}{
			"priorityLevelWithMaxLamports": map[string]interface{}{
				"maxLamports":   300000,
				"priorityLevel": "high",
			},
		}
	}
	raw, err := jup.GetSwapRaw(ctx, swapInput)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(raw)
}

func TestSJup_GetTokenPrice(t *testing.T) {
	ctx := gctx.GetInitCtx()
	jup := NewJup()
	res, err := jup.GetTokenPrice(ctx, "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN,So11111111111111111111111111111111111111112")
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	for k, v := range res.Data {
		g.Dump(k)
		g.Dump(v)
	}
}

func TestSJup_GetPoolInfo(t *testing.T) {
	ctx := gctx.GetInitCtx()
	jup := NewJup()
	res, err := jup.GetPoolInfo(ctx, "2qEHjDLDLbuBgRYvsxhc5D6uDWAivNFZGan56P1tpump")
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(res)
}

func TestSJup_GetTokenInfo(t *testing.T) {
	ctx := gctx.GetInitCtx()
	jup := NewJup()
	res, err := jup.GetTokenInfo(ctx, "So11111111111111111111111111111111111111112")
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(res)
}
