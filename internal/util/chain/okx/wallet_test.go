package okx

import (
	utilCache "keyboard-api-go/internal/util/cache"
	"sync"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	ctx := gctx.New()
	var (
		cache = gcache.New()
		redis = g.Redis()
	)
	// Create redis cache adapter and set it to cache object.
	cache.SetAdapter(gcache.NewAdapterRedis(redis))
	utilCache.InitCache(ctx, cache)
}

func TestOKX_GetWalletAccount(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().GetWalletAccount(ctx, 50, 1)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestOKX_GetPrice(t *testing.T) {
	ctx := gctx.GetInitCtx()
	input := []GetPriceInput{
		{
			ChainIndex: 501,
			Address:    "So11111111111111111111111111111111111111112",
		},
		{
			ChainIndex: 501,
			Address:    "",
		},
	}
	res, err := Okx().GetPrice(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestGetSupportedChains(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().GetSupportedChains(ctx)
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestGetTxHistory(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().TransactionsByAddress(ctx, "501", "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi", "", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestFindTx(t *testing.T) {
	ctx := gctx.GetInitCtx()
	var group sync.WaitGroup
	for i := 0; i < 1; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			res, err := Okx().TransactionDetailByTxHash(ctx, "501",
				"229MHMtFJuxknGFJZzKxj44QkKZ76bkB5w255oACeY3oUbYYeacVcjXRhsRi4JJn9jXW4tf4RkZhoy8Fyityu67w")
			if err != nil {
				t.Fatal(err)
			}
			g.Dump(res)
		}()
	}
	group.Wait()
}

func TestSOkx_TotalValueByAddress(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().TotalValueByAddress(ctx, "501", "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi", "")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestSOkx_AllTokenBalancesByAddress(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().AllTokenBalancesByAddress(ctx, "501", "CjvKgHGWiPq9TCiuDUtDc5LWJydFdLBR7qhH2AzNRZA", "")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestSOkx_TokenBalancesByAddress(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().TokenBalancesByAddress(ctx, "501", "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi", "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestSOkx_TokenDetail(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().TokenDetail(ctx, "501", "So11111111111111111111111111111111111111112")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}

func TestSOkx_TokenList(t *testing.T) {
	ctx := gctx.GetInitCtx()
	res, err := Okx().TokenList(ctx, "501")
	if err != nil {
		t.Fatal(err)
	}
	g.Dump(res)
}
