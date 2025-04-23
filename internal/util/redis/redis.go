package redis

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sRedis struct {
}

type IRedis interface {
	MatchKey(ctx context.Context, key string) []string
}

func New() IRedis {
	return &sRedis{}
}

func (s *sRedis) MatchKey(ctx context.Context, key string) []string {
	redis := g.Redis()
	var cursor int64 = 0
	keys := make([]string, 0)
	for {
		// SCAN 命令进行遍历
		result, err := redis.Do(ctx, "SCAN", cursor, "MATCH", key, "COUNT", 100)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil
		}
		// 解析结果
		data := result.Array()
		cursor = gconv.Int64(data[0])
		keys = append(keys, gconv.Strings(data[1])...)
		if cursor == 0 {
			break
		}
	}
	return keys
}
