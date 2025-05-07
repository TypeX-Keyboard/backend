package ws

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"strings"
	"time"
)

// LoginController 用户登录
func LoginController(client *Client, req *request) {

	userId := gconv.Uint64(0)

	// 用户登录
	login := &login{
		UserId: userId,
		Client: client,
	}
	clientManager.Login <- login

	client.SendMsg(&WResponse{
		Event: Login,
		Data:  "success",
	})
}

func IsAppController(client *Client) {
	client.isApp = true
}

// JoinController 加入
func JoinController(client *Client, req *request) {
	name := gconv.String(req.Data["name"])

	if !client.tags.Contains(name) {
		client.tags.Append(name)
	}
	client.SendMsg(&WResponse{
		Event: Join,
		Data:  client.tags.Slice(),
	})
}

// QuitController 退出
func QuitController(client *Client, req *request) {
	name := gconv.String(req.Data["name"])
	if client.tags.Contains(name) {
		client.tags.RemoveValue(name)
	}
	client.SendMsg(&WResponse{
		Event: Quit,
		Data:  client.tags.Slice(),
	})
}
func PingController(client *Client) {
	currentTime := uint64(gtime.Now().Unix())
	client.Heartbeat(currentTime)
}

func AuthErr(client *Client) {
	client.SendMsg(&WResponse{
		Event: Error,
		Data:  g.Map{"code": 401, "message": "Unauthorized"},
	})
	time.Sleep(5 * time.Second)
	_ = client.Socket.Close()
}

// {"e":"subscribe:mints","d":{"tokens":["So11111111111111111111111111111111111111112"]},"auth":"123"}
func SubscriptionPrice(client *Client, req *request) {
	ctx := context.Background()
	gmlock.Lock(consts.ListenTokenPrices)
	defer gmlock.Unlock(consts.ListenTokenPrices)
	res := WResponse{
		Event: Prices,
	}
	defer func(res *WResponse) {
		client.SendMsg(res)
	}(&res)
	tokens := gconv.Strings(req.Data["tokens"])
	get, err := cache.GetCache().Get(ctx, consts.ListenTokenPrices)
	if err != nil {
		g.Log().Error(ctx, err)
		res.Data = err
		return
	}
	orgTokens := make(map[string]int)
	err = gconv.Struct(get, &orgTokens)
	if err != nil {
		g.Log().Error(ctx, err)
		res.Data = err
		return
	}
	for _, token := range tokens {
		//detail, err := okx.Okx().TokenDetail(ctx, "501", token)
		//if err != nil {
		//	g.Log().Error(ctx, err)
		//	continue
		//}
		//if len(detail.Data) == 0 {
		//	continue
		//}
		if client.Tokens.Contains(token) {
			continue
		}
		if n, ok := orgTokens[token]; ok {
			orgTokens[token] = n + 1
		} else {
			orgTokens[token] = 1
		}
		client.Tokens.Append(token)
	}
	err = cache.GetCache().Set(ctx, consts.ListenTokenPrices, orgTokens, 0)
	if err != nil {
		g.Log().Error(ctx, err)
		res.Data = err
		return
	}
	res.Data = "success"
}

func BroadcastPrices(ctx context.Context) {
	jup := dex.NewJup()
	_, _ = gcron.Add(ctx, "*/5 * * * * *", func(ctx context.Context) {
		get, err := cache.GetCache().Get(ctx, consts.ListenTokenPrices)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		tokens := make(map[string]int)
		err = gconv.Struct(get, &tokens)
		if len(tokens) == 0 {
			return
		}
		keys := make([][]string, 0)
		keySub := make([]string, 0)
		for s, _ := range tokens {
			keySub = append(keySub, s)
			if len(keySub) == 100 {
				keys = append(keys, keySub)
				keySub = make([]string, 0)
			}
		}
		if len(keySub) > 0 {
			keys = append(keys, keySub)
		}
		for _, key := range keys {
			prices, err := jup.GetPoolInfo(ctx, strings.Join(key, ","))
			if err != nil {
				g.Log().Error(ctx, err)
				return
			}
			data := make(map[string]PoolInfo)
			for _, v := range prices.Pools {
				data[v.BaseAsset.ID] = PoolInfo{
					Price:     v.BaseAsset.UsdPrice,
					Sh5m:      v.BaseAsset.Stats5M.PriceChange,
					Sh1h:      v.BaseAsset.Stats1H.PriceChange,
					Sh6h:      v.BaseAsset.Stats6H.PriceChange,
					Sh24h:     v.BaseAsset.Stats5M.PriceChange,
					Mcap:      v.BaseAsset.Mcap,
					Fdv:       v.BaseAsset.Fdv,
					Liquidity: v.Liquidity,
					Volume24H: v.Volume24H,
				}
			}
			SendPrice(&PriceResponse{
				Event: "prices",
				Data:  data,
			})
		}
	})
}
