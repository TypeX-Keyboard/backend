package solana

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"keyboard-api-go/internal/util/chain/solana/model"
	"strings"
	"testing"
)

func init() {

}

func TestDexSwap(t *testing.T) {
	ctx := gctx.New()
	input := model.SwapQuoteReq{
		QuoteBaseReq: model.QuoteBaseReq{
			InputMint:  "So11111111111111111111111111111111111111112",
			OutputMint: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
			Amount:     100000000,
		},
		SlippageBps: 50,
		UsdValue:    1000,
	}
	res, err := NewDex(Jupiter).GetSwapRaw(ctx, input, conf.Address)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
	res2, err := NewDex(Raydium).GetSwapRaw(ctx, input, conf.Address)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res2)
}
func TestDex(t *testing.T) {
	ctx := gctx.New()
	tokens := []string{
		"So11111111111111111111111111111111111111112",
		"FDS6o6yja1ZbRse9UiPTXhG1P7g5qUuSoaT6ciUbJaZw",
		"9BB6NFEcjBCtnNLFko2FqVQBq8HHM13kCyYcdQbgpump",
		"98mb39tPFKQJ4Bif8iVg9mYb9wsfPZgpgN1sxoVTpump",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso8b8JgJV",
		"8PcWAaKGhqP5gfHh6nXTuRxJ9TRLrWVFASZso81b8JgJV",
	}
	res, err := NewDex(Raydium).GetTokenInfo(ctx, strings.Join(tokens, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(res)
	fmt.Println("_________________________________________________________")
	res2, err := NewDex(Jupiter).GetTokenInfo(ctx, strings.Join(tokens, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Dump(res2)
}
