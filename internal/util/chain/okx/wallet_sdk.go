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

// GenerateSignature Generate API request signatures
func (o *sOkx) GenerateSignature(secretKey string, timestamp, method, requestPath string, body []byte) string {
	// 1. Concatenate strings in the prescribed order：timestamp + method + requestPath + body
	message := timestamp + method + requestPath
	if method == "POST" {
		message += string(body)
	}
	//fmt.Println(message)
	// 2. Encrypted with HMAC SHA256
	hmac256 := hmac.New(sha256.New, []byte(secretKey))
	hmac256.Write([]byte(message))

	// 3. Base64 encoding
	signature := base64.StdEncoding.EncodeToString(hmac256.Sum(nil))
	return signature
}

func (o *sOkx) AddHeaders(req *http.Request) (*http.Request, error) {
	// 1. Get the request body
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read the request body: %w", err)
		}
		// Set the request body again, because ReadAll will consume the original request body
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
	// 2. Generate timestamps
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	// 3. Generate signatures
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

// Create an account
func (o *sOkx) CreateWalletAccount(ctx context.Context, accountReq CreateWalletReq) (res CreateWalletRes, err error) {
	res = CreateWalletRes{}
	// Get a proxy
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
	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("The serialization request failed: %v", err))
		return res, err
	}
	req, err := http.NewRequest("POST", "https://www.okx.com/api/v5/wallet/account/create-wallet-account", bytes.NewBuffer(jsonBody))
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("The creation request failed: %v", err))
		return res, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Failed to add request headers: %v", err))
		return res, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("The send request failed: %v", err))
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Read response failed: %v", err))
		return res, err
	}

	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("The API request failed with a status code: %d, response: %s", resp.StatusCode, string(body)))
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Parsing response failed: %v", err))
		return res, err
	}
	return res, nil
}

// Check your account
func (o *sOkx) GetWalletAccount(ctx context.Context, limit int, cursor int) (res GetWalletAccountRes, err error) {
	res = GetWalletAccountRes{}
	// Get a proxy
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
	url := fmt.Sprintf("https://www.okx.com/api/v5/wallet/account/accounts?limit=%d&cursor=%d", limit, cursor)
	//url := fmt.Sprintf("https://www.okx.com/api/v5/wallet/account/accounts")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("The creation request failed: %v", err))
		return res, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Failed to add request headers: %v", err))
		return res, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("The send request failed: %v", err))
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Read response failed: %v", err))
		return res, err
	}

	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("The API request failed with a status code: %d, response: %s", resp.StatusCode, string(body)))
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("Parsing response failed: %v", err))
		return res, err
	}
	return res, nil
}
