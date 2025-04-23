package dex

import (
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/util/chain/solana/model"
	"strconv"
	"testing"
)

func TestRaydiumSwapQuoteEarlyReq(t *testing.T) {
	ctx := gctx.GetInitCtx()
	r := sRaydiumSwapRoute{}
	feeRes, err2 := r.GetPriorityFee(ctx)
	if err2 != nil {
		t.Error(err2)
		return
	}
	g.Dump(feeRes)
	req := model.RaydiumSwapQuoteEarlyReq{
		SwapQuoteReq: model.SwapQuoteReq{
			QuoteBaseReq: model.QuoteBaseReq{
				InputMint:  "So11111111111111111111111111111111111111112",
				OutputMint: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				Amount:     10000,
			},
			SlippageBps: 50,
		},
		TxVersion: "V0",
	}
	EarlyRes, _, err := r.GetSwapQuoteEarly(ctx, req)
	if err != nil {
		t.Error(err)
		return
	}
	//g.Dump(EarlyRes)
	//priceImpactPctStr := strconv.FormatFloat(out.Data.PriceImpactPct, 'f', -1, 64)
	//var rm interface{}
	//err = json.Unmarshal(rrrr, rm)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	for range 10 {
		inp := model.RaydiumSwapQuoteLastReq{
			ComputeUnitPriceMicroLamports: strconv.FormatUint(feeRes.Data.Default.H, 10),
			SwapResponseRaw:               EarlyRes,
			TxVersion:                     "V0",
			Wallet:                        "5dXBG7j4ZHdinmYDLXsPg93B45yv2AHRB5Tqb99GKNiM",
			WrapSol:                       true,
			UnwrapSol:                     true,
		}
		//g.Map{
		//	"": "",
		//}
		marshal, err := json.Marshal(inp)
		if err != nil {
			t.Error(err)
			continue
		}
		g.Log().Info(ctx, string(marshal))
		last, err := r.PostSwapQuoteLast(ctx, inp)
		if err != nil {
			t.Error(err)
			continue
		}
		g.Dump(last)
		break
	}
}

func TestMintsPrice(t *testing.T) {
	ctx := gctx.GetInitCtx()
	r := sRaydiumSwapRoute{}
	tokens := []string{
		"So11111111111111111111111111111111111111112",
		"FDS6o6yja1ZbRse9UiPTXhG1P7g5qUuSoaT6ciUbJaZw",
		"9BB6NFEcjBCtnNLFko2FqVQBq8HHM13kCyYcdQbgpump",
		"98mb39tPFKQJ4Bif8iVg9mYb9wsfPZgpgN1sxoVTpump",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso8b8JgJV",
	}
	price, err := r.MintPrice(ctx, tokens)
	if err != nil {
		t.Error(err)
	}
	g.Dump(price)
}

func TestMintsInfo(t *testing.T) {
	ctx := gctx.GetInitCtx()
	r := sRaydiumSwapRoute{}
	tokens := []string{
		"So11111111111111111111111111111111111111112",
		"FDS6o6yja1ZbRse9UiPTXhG1P7g5qUuSoaT6ciUbJaZw",
		"9BB6NFEcjBCtnNLFko2FqVQBq8HHM13kCyYcdQbgpump",
		"98mb39tPFKQJ4Bif8iVg9mYb9wsfPZgpgN1sxoVTpump",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso8b8JgJV",
	}
	price, err := r.MintInfo(ctx, tokens)
	if err != nil {
		t.Error(err)
	}
	g.Dump(price)
}

func TestRayPrices(t *testing.T) {
	ctx := gctx.GetInitCtx()
	ray := NewRaydium()
	tokens := []string{
		"So11111111111111111111111111111111111111112",
		"FDS6o6yja1ZbRse9UiPTXhG1P7g5qUuSoaT6ciUbJaZw",
		"9BB6NFEcjBCtnNLFko2FqVQBq8HHM13kCyYcdQbgpump",
		"98mb39tPFKQJ4Bif8iVg9mYb9wsfPZgpgN1sxoVTpump",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso8b8JgJV",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso81b8JgJV",
	}
	prices, err := ray.MintPrice(ctx, tokens)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	g.Dump(prices)
}
