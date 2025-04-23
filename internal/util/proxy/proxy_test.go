package proxy

import (
	"context"
	"fmt"
	"keyboard-api-go/internal/consts"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func TestInitProxy(t *testing.T) {
	RedisInit(gctx.GetInitCtx(), UsProxyList, consts.RedisProxyUsKey)
}

func TestProxy(t *testing.T) {
	proxies, err := GetRedisCommonProxies(context.Background(), consts.RedisProxyUsKey)
	fmt.Println(proxies)
	if err != nil {
		return
	}
}

func TestRedisProxyGetByKey(t *testing.T) {
	key, err := RedisProxyGetByKey(gctx.GetInitCtx(), consts.RedisProxyUsKey)
	if err != nil {
		return
	}
	fmt.Println(key)
}
