package okx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"io"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/proxy"
	"net/http"
	"net/url"
	"time"
)

// GenerateSignature 生成 API 请求签名
func (o *sOkx) GenerateSignature(secretKey string, timestamp, method, requestPath string, body []byte) string {
	// 1. 按照规定顺序拼接字符串：timestamp + method + requestPath + body
	message := timestamp + method + requestPath
	if method == "POST" {
		message += string(body)
	}
	//fmt.Println(message)
	// 2. 使用 HMAC SHA256 加密
	hmac256 := hmac.New(sha256.New, []byte(secretKey))
	hmac256.Write([]byte(message))

	// 3. Base64 编码
	signature := base64.StdEncoding.EncodeToString(hmac256.Sum(nil))
	return signature
}

func (o *sOkx) AddHeaders(req *http.Request) (*http.Request, error) {
	// 1. 获取请求体
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("读取请求体失败: %w", err)
		}
		// 重新设置请求体，因为ReadAll会消耗掉原有的请求体
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	for k, v := range o.buildHeaderMap(req.Method, req.URL, bodyBytes) {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (o *sOkx) buildHeaderMap(method string, url *url.URL, bodyBytes []byte) map[string]string {
	if url == nil {
		return nil
	}
	// 2. 生成时间戳
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	// 3. 生成签名
	urlPath := url.Path
	if method == "GET" && len(url.RawQuery) > 0 {
		urlPath = url.Path + "?" + url.RawQuery
	}
	var PopElement = func(slice *[]okxConf) okxConf {
		gmlock.Lock("okx-conf")
		defer gmlock.Unlock("okx-conf")
		if len(*slice) == 0 {
			return okxConf{}
		}
		if len(*slice) == 1 {
			return (*slice)[0]
		}
		var res okxConf
		if len(*slice) > 1 {
			res = (*slice)[0]
			*slice = append((*slice)[1:], (*slice)[0])
		}
		return res
	}
	conf := PopElement(&o.okxConf)
	g.Log().Infof(gctx.New(), "conf: %v", conf)
	signature := o.GenerateSignature(conf.SecretKey, timestamp, method, urlPath, bodyBytes)
	return map[string]string{
		"Content-Type":         "application/json",
		"OK-ACCESS-PROJECT":    conf.ProjectID,
		"OK-ACCESS-KEY":        conf.AccessKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-PASSPHRASE": conf.Passphrase,
		"OK-ACCESS-TIMESTAMP":  timestamp,
	}
}

// 创建账户
func (o *sOkx) CreateWalletAccount(ctx context.Context, accountReq CreateWalletReq) (res CreateWalletRes, err error) {
	res = CreateWalletRes{}
	// 获取代理
	pro, err := proxy.RedisProxyGetByKey(ctx, consts.RedisProxyUsKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return res, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy.GetProxyURL(ctx, pro)),
		},
	}
	requestBody := accountReq
	// 将请求体转换为 JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("序列化请求失败: %v", err))
		return res, err
	}
	req, err := http.NewRequest("POST", "https://www.okx.com/api/v5/wallet/account/create-wallet-account", bytes.NewBuffer(jsonBody))
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("创建请求失败: %v", err))
		return res, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("添加请求头失败: %v", err))
		return res, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("发送请求失败: %v", err))
		return res, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("读取响应失败: %v", err))
		return res, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body)))
		return res, err
	}
	// 解析响应
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v", err))
		return res, err
	}
	return res, nil
}

// 查询账户
func (o *sOkx) GetWalletAccount(ctx context.Context, limit int, cursor int) (res GetWalletAccountRes, err error) {
	res = GetWalletAccountRes{}
	// 获取代理
	//pro, err := proxy.RedisProxyGetByKey(ctx, consts.REDIS_SPIDER_US_PROXY_KEY)
	//if err != nil {
	//	g.Log().Error(ctx, err)
	//	return res, err
	//}
	client := &http.Client{
		Transport: &http.Transport{
			//Proxy: http.ProxyURL(proxy.GetProxyURL(ctx, pro)),
		},
	}
	// 将请求体转换为 JSON
	url := fmt.Sprintf("https://www.okx.com/api/v5/wallet/account/accounts?limit=%d&cursor=%d", limit, cursor)
	//url := fmt.Sprintf("https://www.okx.com/api/v5/wallet/account/accounts")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("创建请求失败: %v", err))
		return res, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("添加请求头失败: %v", err))
		return res, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("发送请求失败: %v", err))
		return res, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("读取响应失败: %v", err))
		return res, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body)))
		return res, err
	}
	// 解析响应
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v", err))
		return res, err
	}
	return res, nil
}
