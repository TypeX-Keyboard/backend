package market

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestQuotes(t *testing.T) {
	ctx := gctx.GetInitCtx()
	slugs := []string{
		"bitcoin",
		"ethereum",
		"solana",
		"bnb",
	}
	r, err := New().Quotes(ctx, slugs)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestMap(t *testing.T) {
	ctx := gctx.GetInitCtx()
	r, err := New().TokenList(ctx, 1, 200)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(len(r))
}
