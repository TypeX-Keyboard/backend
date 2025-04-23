package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/genv"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/service"
	utilCache "keyboard-api-go/internal/util/cache"
	"keyboard-api-go/internal/util/chain"
	"keyboard-api-go/internal/util/proxy"
	"time"
)

func MustInit(ctx context.Context) {
	// 初始化
	initCache(ctx)
	initSetting(ctx)
	initChain(ctx)
	initClientKey(ctx)
	initSave(ctx)
	initBot(ctx)
	proxy.RedisInit(ctx, proxy.UsProxyList, consts.RedisProxyUsKey)
}

func initBot(ctx context.Context) {
	if genv.Get("RE").String() == "" {
		return
	}
	_, err := gcron.Add(ctx, "@hourly", func(ctx context.Context) {
		err := service.Bot().ActiveTask(ctx)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	})
	if err != nil {
		panic(err)
	}
	_, err = gcron.Add(ctx, "0 * * * * *", func(ctx context.Context) {
		err := service.Bot().MiningTask(ctx)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	})
	if err != nil {
		panic(err)
	}
}

// 初始化cache
func initCache(ctx context.Context) {
	var (
		cache = gcache.New()
		redis = g.Redis()
	)
	// Create redis cache adapter and set it to cache object.
	cache.SetAdapter(gcache.NewAdapterRedis(redis))
	utilCache.InitCache(ctx, cache)

	g.Log().Info(ctx, "init cache success")
}

func initClientKey(ctx context.Context) {
	if err := service.Secure().InitClientKey(ctx); err != nil {
		g.Log().Error(ctx, "initClientKey error", err)
	}
	_, err := gcron.Add(ctx, "@daily", func(ctx context.Context) {
		if err := service.Secure().GenClientKey(ctx); err != nil {
			g.Log().Error(ctx, "GenClientKey error", err)
		}
	}, "initClientKey")
	if err != nil {
		panic(err)
	}
}

func initSetting(ctx context.Context) {
	// 初始化设置
	c := g.Config().MustGet(ctx, "setting")
	c.Struct(&model.SettingConfig)
	g.Log().Info(ctx, "init setting success")
}

func initChain(ctx context.Context) {
	c := g.Config().MustGet(ctx, "chain")
	var chainConf chain.ChainConf
	c.Struct(&chainConf)
	if err := chain.InitChain(ctx, chainConf); err != nil {
		panic(err)
	}
	if genv.Get("RE").String() != "" {
		go service.Chain().TxListen(ctx)
		go service.Chain().TxSuccessListen(ctx)
		service.Wallet().Realized(ctx)
		service.Chain().SyncSolPrice(ctx)
	}
}

func initSave(ctx context.Context) {
	re := genv.Get("RE").String()
	if re != "offline" {
		return
	}
	// 立即执行一次
	err := service.ForesightNews().SaveTotalForesightNews(ctx, 1, 10, 100)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 设置每隔 1 小时执行的定时任务
	_, err = gcron.Add(ctx, "*/5 * * * *", func(ctx context.Context) {
		fmt.Println("定时任务执行", time.Now().String())
		// 处理时区锁 0~8：00 无法领取任务
		err = service.ForesightNews().SaveTotalForesightNews(ctx, 1, 1, 100)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		if err != nil {
			g.Log().Error(ctx, err)
		}
	})
	g.Log().Info(ctx, "init task listen success")
}
