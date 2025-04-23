package dex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"
	"keyboard-api-go/internal/util/chain/solana/model"
	"net/http"
	"strings"
)

type SwapRoute interface {
	GetSwapRoute(ctx context.Context, input interface{}) (interface{}, error)
}

const (
	BASE_HOST           = "https://api-v3.raydium.io"
	OWNER_BASE_HOST     = "https://owner-v1.raydium.io"
	SERVICE_BASE_HOST   = "https://service.raydium.io"
	MONITOR_BASE_HOST   = "https://monitor.raydium.io"
	SERVICE_1_BASE_HOST = "https://service-v1.raydium.io"

	SEND_TRANSACTION = "/send-transaction"
	FARM_ARP         = "/main/farm/info"
	FARM_ARP_LINE    = "/main/farm-apr-tv"

	CLMM_CONFIG = "/main/clmm-config"
	CPMM_CONFIG = "/main/cpmm-config"

	VERSION = "/main/version"

	// api v3
	CHECK_AVAILABILITY = "/v3/main/AvailabilityCheckAPI"
	RPCS               = "/main/rpcs"
	INFO               = "/main/info"
	STAKE_POOLS        = "/main/stake-pools"
	CHAIN_TIME         = "/main/chain-time"

	TOKEN_LIST   = "/mint/list"
	MINT_INFO_ID = "/mint/ids"

	JUP_TOKEN_LIST = "https://tokens.jup.ag/tokens?tags=lstcommunity"
	/**
	 * poolType =  {all concentrated standard allFarm concentratedFarm standardFarm}
	 * poolSortField =  {liquidity | volume_24h / 7d / 30d | fee_24h / 7d / 30d | apr_24h / 7d / 30d}
	 * sortType =  {desc/asc}
	 * page =  number
	 * pageSize =  number
	 */
	POOL_LIST = "/pools/info/list"
	/**
	 * ?ids=idList.join('')
	 */
	POOL_SEARCH_BY_ID = "/pools/info/ids"
	/**
	 * mint1/mint2 =  search pool by mint
	 * poolSortField =  {liquidity | volume_24h / 7d / 30d | fee_24h / 7d / 30d | apr_24h / 7d / 30d}
	 * poolType =  {all concentrated standard allFarm concentratedFarm standardFarm}
	 * sortType =  {desc/asc}
	 * page =  number
	 * pageSize =  number
	 */
	POOL_SEARCH_MINT = "/pools/info/mint"
	/** ?lps=lpList.join('') */
	POOL_SEARCH_LP = "/pools/info/lps"
	/** ?ids=idList.join(',') */
	POOL_KEY_BY_ID = "/pools/key/ids"
	/** ?id=string */
	POOL_LIQUIDITY_LINE = "/pools/line/liquidity"
	POOL_POSITION_LINE  = "/pools/line/position"

	FARM_INFO = "/farms/info/ids"
	/** ?lp=string&pageSize=100&page=number */
	FARM_LP_INFO = "/farms/info/lp"
	FARM_KEYS    = "/farms/key/ids"

	OWNER_CREATED_FARM  = "/create-pool/{owner}"
	OWNER_IDO           = "/main/ido/{owner}"
	OWNER_STAKE_FARMS   = "/position/stake/{owner}"
	OWNER_LOCK_POSITION = "/position/clmm-lock/{owner}"
	IDO_KEYS            = "/ido/key/ids"
	SWAP_HOST           = "https://transaction-v1.raydium.io"
	SWAP_COMPUTE        = "/compute/"
	SWAP_TX             = "/transaction/"
	MINT_PRICE          = "/mint/price"
	MIGRATE_CONFIG      = "/main/migrate-lp"
	PRIORITY_FEE        = "/main/auto-fee"

	CPMM_LOCK = "https://dynamic-ipfs.raydium.io/lock/cpmm/position"
)

type sRaydiumSwapRoute struct {
}

type IRaydiumSwapRoute interface {
	GetPriorityFee(ctx context.Context) (model.RaydiumPriorityFeeBaseRes, error)
	GetSwapQuoteEarly(ctx context.Context, input model.RaydiumSwapQuoteEarlyReq) (model.RaydiumSwapQuoteEarlyRes, []byte, error)
	PostSwapQuoteLast(ctx context.Context, input model.RaydiumSwapQuoteLastReq) (model.RaydiumSwapQuoteLastRes, error)
	MintPrice(ctx context.Context, tokenAddress []string) (map[string]float64, error)
	MintInfo(ctx context.Context, tokenAddress []string) ([]model.RaydiumMintInfoData, error)
}

func NewRaydium() IRaydiumSwapRoute {
	return &sRaydiumSwapRoute{}
}

func (r *sRaydiumSwapRoute) GetPriorityFee(ctx context.Context) (model.RaydiumPriorityFeeBaseRes, error) {
	res := model.RaydiumPriorityFeeBaseRes{}
	url := BASE_HOST + PRIORITY_FEE
	var response *gclient.Response
	var err error
	for range 5 {
		response, err = g.Client().Get(ctx, url, nil)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return res, err
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return res, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()
	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	return res, err
}

func (r *sRaydiumSwapRoute) GetSwapQuoteEarly(ctx context.Context, input model.RaydiumSwapQuoteEarlyReq) (model.RaydiumSwapQuoteEarlyRes, []byte, error) {
	res := model.RaydiumSwapQuoteEarlyRes{}
	urlFormat := SWAP_HOST + "/compute/swap-base-in?inputMint=%s&outputMint=%s&amount=%d&slippageBps=%d&txVersion=%s"
	// 使用初始化后的 swapBaseUrlFormat 字段构建请求 URL
	url := fmt.Sprintf(urlFormat, input.InputMint, input.OutputMint, input.Amount, input.SlippageBps, input.TxVersion)
	g.Log().Info(ctx, url)
	client := g.Client()
	var response *gclient.Response
	var err error
	for range 5 {
		response, err = client.Get(ctx, url)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return res, nil, err
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return res, nil, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()

	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	return res, responseBody, nil
}

func (r *sRaydiumSwapRoute) PostSwapQuoteLast(ctx context.Context, input model.RaydiumSwapQuoteLastReq) (model.RaydiumSwapQuoteLastRes, error) {
	res := model.RaydiumSwapQuoteLastRes{}
	url := SWAP_HOST + "/transaction/swap-base-in"
	g.Log().Info(ctx, url)
	client := g.Client()
	marshal, _ := json.Marshal(input)
	var response *gclient.Response
	var err error
	for range 5 {
		response, err = client.Post(ctx, url, marshal)
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
		g.Log().Info(ctx, string(responseBody))
		return res, gerror.New("http status code not 200, " + response.Status)
	}
	err = gjson.Unmarshal(responseBody, &res)
	return res, err
}

func (r *sRaydiumSwapRoute) MintPrice(ctx context.Context, tokenAddress []string) (map[string]float64, error) {
	res := model.RaydiumMintPriceRes{}
	url := fmt.Sprintf("%s?mints=%s", BASE_HOST+MINT_PRICE, strings.Join(tokenAddress, ","))
	var response *gclient.Response
	var err error
	for range 5 {
		response, err = g.Client().Get(ctx, url)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return nil, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()
	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	data := make(map[string]float64)
	for mint, v := range res.Data {
		data[mint] = gconv.Float64(v)
	}
	return data, err
}

func (r *sRaydiumSwapRoute) MintInfo(ctx context.Context, tokenAddress []string) ([]model.RaydiumMintInfoData, error) {
	res := model.RaydiumMintInfoRes{}
	url := fmt.Sprintf("%s?mints=%s", BASE_HOST+MINT_INFO_ID, strings.Join(tokenAddress, ","))
	var response *gclient.Response
	var err error
	for range 5 {
		response, err = g.Client().Get(ctx, url)
		if err == nil {
			break
		}
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	defer response.Close()
	if response.StatusCode != http.StatusOK {
		return nil, gerror.New("http status code not 200, " + response.Status)
	}
	responseBody := response.ReadAll()
	//g.Log().Info(ctx, string(responseBody))
	err = gjson.Unmarshal(responseBody, &res)
	for i, datum := range res.Data {
		if len(datum.Address) == 0 {
			res.Data = append(res.Data[:i], res.Data[i+1:]...)
		}
	}
	return res.Data, err
}
