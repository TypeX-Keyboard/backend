package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"keyboard-api-go/internal/util/cache"
	"net/url"
	"os"
)

type Proxy struct {
	Server   string
	Username string
	Password string
	Http     string
	Https    string
	Proxies  map[string]string
}

func (p *Proxy) SetProxies() {
	p.Http = fmt.Sprintf("http://%s:%s@%s", p.Username, p.Password, p.Server)
	p.Https = fmt.Sprintf("https://%s:%s@%s", p.Username, p.Password, p.Server)
	p.Proxies = map[string]string{
		"http":  p.Http,
		"https": p.Https,
	}
}

// 美国 amazon、tiktok、youtube通杀代理
var UsProxyList = []Proxy{
	// add your proxy
}

// 马来 吉隆坡  1g流量
var myProxyList = []Proxy{}

// 8.7 日本 神奈川  1g流量
var jpProxyList = []Proxy{}

// 新加坡
var sgProxyList = []Proxy{}

var totalProxyList = [][]Proxy{UsProxyList, sgProxyList, myProxyList, jpProxyList}

func RedisInit(ctx context.Context, proxyList []Proxy, key string) {
	redis := g.Redis()
	// 删除已经存在的键
	_, err := redis.Del(ctx, key)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	for _, proxy := range proxyList {
		// 如果已存在 则更新
		temp := gredis.ZAddMember{Member: proxy, Score: float64(gtime.Now().Timestamp())}
		_, err := redis.ZAdd(ctx, key, nil, temp)
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}
}

// GetRedisCommonProxies 从redis中获取一个代理服务器
func GetRedisCommonProxies(ctx context.Context, key string) (Proxy, error) {
	res, err := cache.GetCache().Get(ctx, key)
	if err != nil {
		g.Log().Error(ctx, err)
		return Proxy{}, err
	}
	var proxy Proxy
	err = res.Struct(&proxy)
	if err != nil {
		g.Log().Error(ctx, err)
		return Proxy{}, err
	}
	return proxy, nil
}

// RedisProxyGetByKey 从redis中根据key获取一个代理服务器
func RedisProxyGetByKey(ctx context.Context, key string) (Proxy, error) {
	//return Proxy{}, errors.New("没有可用的代理")
	redis := g.Redis()
	revRange, err := redis.ZRange(ctx, key, 0, 0) // 取低分
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if len(revRange) == 0 {
		return Proxy{}, errors.New("没有可用的代理")
	}
	// 将返回的数据转换为 []interface{} 类型
	var p Proxy
	for _, val := range revRange {
		err = val.Struct(&p)
		if err != nil {
			g.Log().Error(ctx, err)

		}
		f := float64(1)
		// fmt.Printf("代理IP 为 %s时间戳：%f\n", p.Server, f)
		_, err = redis.ZIncrBy(ctx, key, f, p)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		p.SetProxies()
		break
	}
	return p, nil
}

func GetProxyURL(ctx context.Context, proxy Proxy) *url.URL {
	re := os.Getenv("RE")
	proxyURL, _ := url.Parse(proxy.Http)
	//g.Log().Debug(ctx, "代理环境 为 ", re)
	if re == "" {
		proxyURL, _ = url.Parse("http://127.0.0.1:7890")
	}
	return proxyURL
}
