package dex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/chain/solana/model"
	"keyboard-api-go/internal/util/proxy"
	"net/http"
)

const MAX_RETRIES = 5

type sJup struct {
}

type IJup interface {
	GetSwapQuoteEarly(ctx context.Context, input model.SwapQuoteReq) (res model.JupSwapQuoteRes, err error)
	GetSwapRaw(ctx context.Context, input model.GetSwapRawReq) (res model.GetSwapRawRes, err error)
	GetTokenPrice(ctx context.Context, tokens string) (res model.JupiterPriceV2Res, err error)
	GetPoolInfo(ctx context.Context, tokens string) (*model.GetPoolInfoRes, error)
	GetTokenInfo(ctx context.Context, token string) (*model.JupTokenInfo, error)
}

func NewJup() IJup {
	return &sJup{}
}

func (j *sJup) GetSwapQuoteEarly(ctx context.Context, input model.SwapQuoteReq) (res model.JupSwapQuoteRes, err error) {
	client := g.Client()
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err == nil {
		client.Proxy(pro.Http)
	}
	var response *gclient.Response
	payload := g.Map{
		"inputMint":      input.InputMint,
		"outputMint":     input.OutputMint,
		"amount":         input.Amount,
		"autoSlippage":   true,
		"platformFeeBps": 50,
	}
	if input.UsdValue > 0 {
		payload["autoSlippageCollisionUsdValue"] = input.UsdValue
	}
	for range 5 {
		response, err = client.Get(ctx, "https://api.jup.ag/swap/v1/quote", payload)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return res, err
	}
	defer response.Close()
	responseBody := response.ReadAll()
	if response.StatusCode != http.StatusOK {
		g.Log().Error(ctx, string(responseBody))
		return res, gerror.New("http status code not 200, " + response.Status)
	}

	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	return res, nil
}

func (j *sJup) GetSwapRaw(ctx context.Context, input model.GetSwapRawReq) (res model.GetSwapRawRes, err error) {
	url := "https://api.jup.ag/swap/v1/swap"
	marshal, _ := json.Marshal(input)
	client := g.Client()
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err == nil {
		client.Proxy(pro.Http)
	}
	var response *gclient.Response
	for range 5 {
		response, err = client.Post(ctx, url, marshal)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return res, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()

	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	return
}

func (j *sJup) GetTokenPrice(ctx context.Context, tokens string) (res model.JupiterPriceV2Res, err error) {
	res = model.JupiterPriceV2Res{}
	url := "https://api.jup.ag/price/v2"
	params := g.Map{
		"ids":           tokens,
		"showExtraInfo": true,
	}
	client := g.Client()
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err == nil {
		client.Proxy(pro.Http)
	}
	var finalErr error
	for range MAX_RETRIES {
		response, err := client.Get(ctx, url, params)
		if err != nil {
			g.Log().Error(ctx, err)
			finalErr = err
			continue
		}
		// Use anonymous functions to ensure that the response is closed at the end of the current loop
		func() {
			defer response.Close()
			if response.StatusCode != http.StatusOK {
				finalErr = fmt.Errorf("http status code not 200: %s", response.Status)
				g.Log().Error(ctx, finalErr)
				return
			}
			responseBody := response.ReadAll()
			//g.Log().Info(ctx, string(responseBody))
			if err := gjson.Unmarshal(responseBody, &res); err != nil {
				finalErr = err
				g.Log().Error(ctx, err)
				return
			}
			finalErr = nil // Clear the error on success
		}()
		if finalErr == nil {
			break // 成功则退出循环
		}
	}
	return res, finalErr
}
func (j *sJup) GetPoolInfo(ctx context.Context, tokens string) (*model.GetPoolInfoRes, error) {
	headerMap := map[string]string{
		"accept":             "application/json",
		"accept-language":    "zh-CN,zh;q=0.9",
		"origin":             "https://jup.ag",
		"priority":           "u=1, i",
		"referer":            "https://jup.ag/",
		"sec-ch-ua":          `Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	}
	c := g.Client()
	c.SetHeaderMap(headerMap)
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err == nil {
		c.Proxy(pro.Http)
	}
	response, err := c.Get(ctx, "https://datapi.jup.ag/v1/pools", g.Map{
		"assetIds": tokens,
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return nil, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()
	var res model.GetPoolInfoRes
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return &res, nil
}

func (j *sJup) GetTokenInfo(ctx context.Context, token string) (*model.JupTokenInfo, error) {
	c := g.Client()
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err == nil {
		c.Proxy(pro.Http)
	}
	get, err := c.Get(ctx, fmt.Sprintf("https://api.jup.ag/tokens/v1/token/%s", token))
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		return nil, gerror.New("http status code not 200, " + get.Status)
	}
	responseBody := get.ReadAll()
	var res model.JupTokenInfo
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return &res, nil
}
