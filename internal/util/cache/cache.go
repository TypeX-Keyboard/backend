package cache

import (
	"context"
	"fmt"
	"keyboard-api-go/internal/consts"

	"github.com/gogf/gf/v2/os/gcache"
)

var cache *gcache.Cache

func InitCache(ctx context.Context, adapter *gcache.Cache) {
	cache = adapter
}

func GetCache() *gcache.Cache {
	return cache
}

func GetAesKeyCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", consts.AesKeyCacheKeyPrefix, uuid)
}
