package secure

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/secure"
	"time"
)

func init() {
	service.RegisterSecure(New())
}

func New() service.ISecure {
	return &sSecure{}
}

type sSecure struct {
}

func (s *sSecure) AcquirePublicKey(ctx context.Context) (string, error) {
	if model.SettingConfig == nil {
		return "", fmt.Errorf("model.SettingConfig is nil")
	}
	return consts.PublicKey, nil
}

func (s *sSecure) SubmitKey(ctx context.Context, uuid, key string) error {
	g.Log().Infof(ctx, "SubmitKey: %s, %s", uuid, key)
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	aesKey, err := secure.Sha256Decrypt(data, []byte(consts.PrivateKey))
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if aesKey == nil || len(aesKey) != 32 {
		return fmt.Errorf("invalid key")
	}
	err = cache.GetCache().Set(ctx, cache.GetAesKeyCacheKey(uuid), aesKey, 30*24*time.Hour)
	return nil
}

func (s *sSecure) GenClientKey(ctx context.Context) error {
	privateKeyB, publicKeyB, err := secure.GenerateKeyPair(2048)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyB)
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyB)
	if err := cache.GetCache().Set(ctx, consts.RedisClientPrivateKey, privateKeyStr, time.Hour*consts.ExpireH); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if err := cache.GetCache().Set(ctx, consts.RedisClientPublicKey, publicKeyStr, time.Hour*consts.ExpireH); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	key, err := secure.GeneratedSecretKey(32)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if err := cache.GetCache().Set(ctx, consts.RedisClientHMACKey, key, time.Hour*consts.ExpireH); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	return nil
}

func (s *sSecure) InitClientKey(ctx context.Context) error {
	get, err := cache.GetCache().Get(ctx, consts.RedisClientPrivateKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	if !get.IsNil() {
		return nil
	}
	return s.GenClientKey(ctx)
}
