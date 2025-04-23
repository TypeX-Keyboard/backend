package okx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"keyboard-api-go/internal/util/str"
	"net/http"
	"net/url"
	"strings"
)

func (o *sOkx) GetPrice(ctx context.Context, input []GetPriceInput) (*PriceRes, error) {
	// 获取代理
	//pro, err := proxy.RedisProxyGetByKey(ctx, consts.REDIS_SPIDER_US_PROXY_KEY)
	//if err != nil {
	//	g.Log().Error(ctx, err)
	//	return nil, err
	//}
	client := &http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(proxy.GetProxyURL(ctx, pro)),
		//},
	}
	// 将请求体转换为 JSON
	jsonBody, err := json.Marshal(input)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("序列化请求失败: %v", err))
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.okx.com/api/v5/wallet/token/real-time-price", bytes.NewBuffer(jsonBody))
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("创建请求失败: %v", err))
		return nil, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("添加请求头失败: %v", err))
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("发送请求失败: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("读取响应失败: %v", err))
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body)))
		return nil, err
	}
	// 解析响应
	res := PriceRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, string(body)))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) GetSupportedChains(ctx context.Context) ([]Chain, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.okx.com/api/v5/wallet/chain/supported-chains", nil)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("创建请求失败: %v", err))
		return nil, err
	}
	req, err = o.AddHeaders(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("添加请求头失败: %v", err))
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("发送请求失败: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("读取响应失败: %v", err))
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body)))
		return nil, err
	}
	// 解析响应
	res := ChainRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, string(body)))
		return nil, err
	}
	return res.Data, nil
}

func (o *sOkx) TotalValueByAddress(ctx context.Context, chain, address, assetType string, excludeRiskToken ...bool) (*TotalValueByAddressRes, error) {
	params := g.Map{
		"address": address,
		"chains":  chain,
	}
	if len(assetType) > 0 {
		params["assetType"] = assetType
	}
	if len(excludeRiskToken) > 0 {
		params["excludeRiskToken"] = excludeRiskToken[0]
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/asset/total-value-by-address", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res TotalValueByAddressRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) AllTokenBalancesByAddress(ctx context.Context, chain, address, filter string) (*AllTokenBalancesByAddressRes, error) {
	params := g.Map{
		"address": address,
		"chains":  chain,
	}
	if len(filter) > 0 {
		params["filter"] = filter
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/asset/all-token-balances-by-address", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res AllTokenBalancesByAddressRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) TokenBalancesByAddress(ctx context.Context, chain, address, tokenAddresses string) (*AllTokenBalancesByAddressRes, error) {
	tokens := make([]g.Map, 0)
	for _, s := range strings.Split(tokenAddresses, ",") {
		tokens = append(tokens, g.Map{
			"chainIndex":   chain,
			"tokenAddress": s,
		})
	}
	params := g.Map{
		"address":        address,
		"tokenAddresses": tokens,
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/asset/token-balances-by-address", nil))
	if err != nil {
		return nil, err
	}
	marshal, _ := json.Marshal(params)
	c.SetHeaderMap(o.buildHeaderMap("POST", link, marshal))
	get, err := c.Post(ctx, link.String(), params)
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res AllTokenBalancesByAddressRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) TransactionsByAddress(ctx context.Context, chain, address, tokenAddress, begin, end, cursor, limit string) (*TxHistoryRes, error) {
	params := g.Map{
		"address": address,
		"chains":  chain,
	}
	if len(tokenAddress) > 0 {
		params["tokenAddress"] = tokenAddress
	}
	if len(begin) > 0 {
		params["begin"] = begin
	}
	if len(end) > 0 {
		params["end"] = end
	}
	if len(cursor) > 0 {
		params["cursor"] = cursor
	}
	if len(limit) > 0 {
		params["limit"] = limit
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/post-transaction/transactions-by-address", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res TxHistoryRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) TransactionDetailByTxHash(ctx context.Context, chain, txHash string) (*TxHashRes, error) {
	params := g.Map{
		"txHash":     txHash,
		"chainIndex": chain,
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/post-transaction/transaction-detail-by-txhash", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res TxHashRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) TokenDetail(ctx context.Context, chainIndex, tokenAddress string) (*TokenDetailRes, error) {
	params := g.Map{
		"tokenAddress": tokenAddress,
		"chainIndex":   chainIndex,
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/wallet/token/token-detail", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res TokenDetailRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}

func (o *sOkx) TokenList(ctx context.Context, chainId string) (*TokenListRes, error) {
	params := g.Map{
		"chainId": chainId,
	}
	c := g.Client()
	link, err := url.Parse(str.AddQuery("https://www.okx.com/api/v5/dex/aggregator/all-tokens", params))
	if err != nil {
		return nil, err
	}

	c.SetHeaderMap(o.buildHeaderMap("GET", link, nil))
	get, err := c.Get(ctx, link.String())
	if err != nil {
		return nil, err
	}
	defer get.Close()
	if get.StatusCode != http.StatusOK {
		g.Log().Error(ctx, fmt.Sprintf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString()))
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", get.StatusCode, get.ReadAllString())
	}
	var res TokenListRes
	err = json.Unmarshal(get.ReadAll(), &res)
	if err != nil {
		g.Log().Error(ctx, fmt.Sprintf("解析响应失败: %v, 响应: %s", err, get.ReadAllString()))
		return nil, err
	}
	return &res, nil
}
