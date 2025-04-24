package secure

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/cache"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"keyboard-api-go/api/secure/v1"
)

func (c *ControllerV1) GetKey(ctx context.Context, req *v1.GetKeyReq) (res *v1.GetKeyRes, err error) {
	res = &v1.GetKeyRes{}
	var redisKey string
	if req.From == "native" {
		redisKey = consts.RedisClientPrivateKey
	} else if req.From == "web" {
		redisKey = consts.RedisClientPublicKey
	} else {
		return res, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	get, err := cache.GetCache().Get(ctx, redisKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	if get == nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.Key = get.String()
	get, err = cache.GetCache().Get(ctx, consts.RedisClientHMACKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	if get == nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.SignKey = get.String()
	expire, err := cache.GetCache().GetExpire(ctx, redisKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, gerror.NewCode(gcode.CodeInternalError)
	}
	res.Expire = time.Now().Add(expire).Unix()
	return res, gerror.NewCode(gcode.CodeOK)
}
