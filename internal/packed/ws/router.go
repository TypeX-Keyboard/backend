package ws

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/cache"
)

const (
	Error  = "error"
	Login  = "login"
	Join   = "join"
	Quit   = "quit"
	IsApp  = "is_app"
	Ping   = "ping"
	Prices = "subscribe:mints"
)

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()
	request := &request{}
	err := gconv.Struct(message, request)
	if err != nil {
		fmt.Println("数据解析失败：", err)
		return
	}
	client.AuthToken = request.AuthToken
	ctx := gctx.New()
	get, err := cache.GetCache().Get(ctx, consts.RedisClientHMACKey)
	if err != nil {
		g.Log().Error(ctx, err)
		go AuthErr(client)
		return
	}
	if get == nil {
		g.Log().Error(ctx, err)
		go AuthErr(client)
		return
	}
	// 2. 验证 Token
	if client.AuthToken != get.String() {
		go AuthErr(client)
		return
	}
	switch request.Event {
	case Login:
		LoginController(client, request)
		break
	case Join:
		JoinController(client, request)
		break
	case Quit:
		QuitController(client, request)
		break
	case IsApp:
		IsAppController(client)
		break
	case Ping:
		PingController(client)
		break
	case Prices:
		SubscriptionPrice(client, request)
		break
	}
}
