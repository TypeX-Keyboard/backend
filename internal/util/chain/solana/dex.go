package solana

import (
	"context"
	"errors"
	"fmt"
	"keyboard-api-go/internal/util/chain/solana/dex"
	"keyboard-api-go/internal/util/chain/solana/model"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shopspring/decimal"
)

type IDex interface {
	GetPriorityFee(ctx context.Context) (model.PriorityFeeRes, error)
	GetSlippageBps(ctx context.Context, input model.QuoteBaseReq) (uint64, error)
	GetSwapRaw(ctx context.Context, input model.SwapQuoteReq, wallet string) (res model.GetSwapRouterRes, err error)
	GetTokenPrice(ctx context.Context, tokens string) (res map[string]float64, err error)
	GetTokenInfo(ctx context.Context, tokens string) (res []model.TokenInfoRes, err error)
}

type sDex struct {
	adaptor Adaptor
}

type Adaptor string

var Jupiter = Adaptor("Jupiter")
var Raydium = Adaptor("Raydium")

func NewDex(adaptor Adaptor) IDex {
	return &sDex{
		adaptor,
	}
}

func (s *sDex) GetPriorityFee(ctx context.Context) (res model.PriorityFeeRes, err error) {
	r := dex.NewRaydium()
	feeRes, err := r.GetPriorityFee(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return feeRes.Data.Default, nil
}

func (s *sDex) GetSlippageBps(ctx context.Context, input model.QuoteBaseReq) (uint64, error) {
	jup := dex.NewJup()
	out, err := GetSolana().GetTokenSupply(ctx, input.InputMint)
	if err != nil {
		return 0, err
	}
	prices, err := s.GetTokenPrice(ctx, input.InputMint)
	if err != nil {
		return 0, err
	}
	price := decimal.NewFromFloat(prices[input.InputMint])
	amount := decimal.NewFromInt(int64(input.Amount)).Shift(-int32(out.Value.Decimals))
	usdValue := price.Mul(amount)
	var usd uint64 = 1
	if usdValue.IntPart() > 0 {
		usd = uint64(usdValue.IntPart())
	}
	quote, err := jup.GetSwapQuoteEarly(ctx, model.SwapQuoteReq{
		QuoteBaseReq: input,
		UsdValue:     usd,
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, err
	}
	if quote.ComputedAutoSlippage == 0 {
		quote.ComputedAutoSlippage = quote.SlippageBps
	}
	return quote.ComputedAutoSlippage, nil
}

func (s *sDex) GetSwapRaw(ctx context.Context, input model.SwapQuoteReq, wallet string) (res model.GetSwapRouterRes, err error) {
	if input.InputMint == input.OutputMint {
		return res, errors.New("input and output token is same")
	}
	switch s.adaptor {
	case Jupiter:
		return s.jupSwapRaw(ctx, input, wallet)
	case Raydium:
		return s.raySwapRaw(ctx, input, wallet)
	default:
		return res, fmt.Errorf("not support adaptor %s", s.adaptor)
	}
}

func (s *sDex) jupSwapRaw(ctx context.Context, input model.SwapQuoteReq, wallet string) (res model.GetSwapRouterRes, err error) {
	res.Wallet = wallet
	res.InputMint = input.InputMint
	res.OutputMint = input.OutputMint
	res.PriorityFee = input.PriorityFee
	jup := dex.NewJup()
	quote, err := jup.GetSwapQuoteEarly(ctx, model.SwapQuoteReq{
		QuoteBaseReq: model.QuoteBaseReq{
			InputMint:  input.InputMint,
			OutputMint: input.OutputMint,
			Amount:     input.Amount,
		},
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res.InAmount = quote.InAmount
	res.OutAmount = quote.OutAmount
	res.SlippageBps = input.SlippageBps
	res.RoutePlan = make([]model.RoutePlan, 0)
	for _, routePlan := range quote.RoutePlan {
		res.RoutePlan = append(res.RoutePlan, model.RoutePlan{
			InputMint:  routePlan.SwapInfo.InputMint,
			OutputMint: routePlan.SwapInfo.OutputMint,
			FeeMint:    routePlan.SwapInfo.FeeMint,
			FeeAmount:  routePlan.SwapInfo.FeeAmount,
		})
	}
	swapInput := model.GetSwapRawReq{
		QuoteResponse:           quote,
		FeeAccount:              input.FeeAccount,
		UserPublicKey:           wallet,
		UseSharedAccounts:       true,
		DynamicComputeUnitLimit: true,
		WrapAndUnwrapSol:        true,
		DynamicSlippage:         true,
	}
	if input.Mev {
		swapInput.PrioritizationFeeLamports = map[string]interface{}{
			"jitoTipLamports": input.PriorityFee,
		}
	} else {
		swapInput.PrioritizationFeeLamports = map[string]interface{}{
			"priorityLevelWithMaxLamports": map[string]interface{}{
				"maxLamports":   input.PriorityFee,
				"priorityLevel": "high",
			},
		}
	}
	raw, err := jup.GetSwapRaw(ctx, swapInput)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res.TxRaw = raw.SwapTransaction
	return res, nil
}

func (s *sDex) raySwapRaw(ctx context.Context, input model.SwapQuoteReq, wallet string) (res model.GetSwapRouterRes, err error) {
	res.Wallet = wallet
	res.InputMint = input.InputMint
	res.OutputMint = input.OutputMint
	res.PriorityFee = input.PriorityFee
	r := dex.NewRaydium()
	req := model.RaydiumSwapQuoteEarlyReq{
		SwapQuoteReq: model.SwapQuoteReq{
			QuoteBaseReq: model.QuoteBaseReq{
				InputMint:  input.InputMint,
				OutputMint: input.OutputMint,
				Amount:     input.Amount,
			},
			SlippageBps: input.SlippageBps,
			PriorityFee: input.PriorityFee,
		},
		TxVersion: "V0",
	}
	quote, _, err := r.GetSwapQuoteEarly(ctx, req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res.InAmount = quote.Data.InputAmount
	res.OutAmount = quote.Data.OutputAmount
	res.SlippageBps = input.SlippageBps
	res.RoutePlan = make([]model.RoutePlan, 0)
	for _, routePlan := range quote.Data.RoutePlan {
		res.RoutePlan = append(res.RoutePlan, model.RoutePlan{
			InputMint:  routePlan.InputMint,
			OutputMint: routePlan.OutputMint,
			FeeMint:    routePlan.FeeMint,
			FeeAmount:  routePlan.FeeAmount,
		})
	}
	inp := model.RaydiumSwapQuoteLastReq{
		ComputeUnitPriceMicroLamports: strconv.FormatUint(input.PriorityFee, 10),
		SwapResponseRaw:               quote,
		TxVersion:                     "V0",
		Wallet:                        wallet,
		WrapSol:                       true,
		UnwrapSol:                     true,
	}
	last, err := r.PostSwapQuoteLast(ctx, inp)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if len(last.Data) == 0 {
		err = errors.New("no tx found")
		g.Log().Error(ctx, err)
		return
	}
	res.TxRaw = last.Data[0].Transaction
	return res, nil
}

func (s *sDex) GetTokenInfo(ctx context.Context, tokens string) (res []model.TokenInfoRes, err error) {
	ray := dex.NewRaydium()
	info, err := ray.MintInfo(ctx, strings.Split(tokens, ","))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	prices, err := s.GetTokenPrice(ctx, tokens)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res = make([]model.TokenInfoRes, 0)
	for _, v := range info {
		res = append(res, model.TokenInfoRes{
			RaydiumMintInfoData: v,
			Price:               prices[v.Address],
		})
	}
	return
}

func (s *sDex) GetTokenPrice(ctx context.Context, tokens string) (res map[string]float64, err error) {
	switch s.adaptor {
	case Jupiter:
		return s.jupTokenPrice(ctx, tokens)
	case Raydium:
		return s.rayTokenPrice(ctx, tokens)
	default:
		return res, fmt.Errorf("not support adaptor %s", s.adaptor)
	}
}

func (s *sDex) jupTokenPrice(ctx context.Context, tokens string) (map[string]float64, error) {
	jup := dex.NewJup()
	jupRes, err := jup.GetTokenPrice(ctx, tokens)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	res := make(map[string]float64)
	for mint, v := range jupRes.Data {
		res[mint] = gconv.Float64(v.Price)
	}
	return res, err
}

func (s *sDex) rayTokenPrice(ctx context.Context, tokens string) (map[string]float64, error) {
	ray := dex.NewRaydium()
	split := strings.Split(tokens, ",")
	rayRes, err := ray.MintPrice(ctx, split)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return rayRes, err
}
