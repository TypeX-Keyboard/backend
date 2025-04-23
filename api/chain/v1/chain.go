package v1

import (
	"keyboard-api-go/internal/model"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gogf/gf/v2/frame/g"
)

type SolTransferReq struct {
	g.Meta `path:"/transfer" tags:"区块链" method:"post" summary:"Sol 转账Tx"`
	From   string `json:"from" dc:"发送地址" required:"true"`
	To     string `json:"to" dc:"接收地址" required:"true"`
	Token  string `json:"token" dc:"代币地址 sol地址(So11111111111111111111111111111111111111112)" required:"true"`
	Amount uint64 `json:"amount" dc:"转账数量，需要乘token的decimal" required:"true"`
}

type SolTransferRes struct {
	g.Meta `mime:"application/json"`
	model.TransferRes
}

type SolPriorityFeeReq struct {
	g.Meta `path:"/solPriorityFee" tags:"区块链" method:"get" summary:"Sol 获取优先交易费"`
}

type SolPriorityFeeRes struct {
	g.Meta `mime:"application/json"`
	model.PriorityFeeRes
}

type SolSlippageBpsReq struct {
	g.Meta   `path:"/solSlippageBps" tags:"区块链" method:"get" summary:"Sol 获取动态滑点"`
	TokenIn  string `json:"token_in" dc:"要支付的代币地址 sol地址(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut string `json:"token_out" dc:"要获得的代币地址" required:"true"`
	AmountIn uint64 `json:"amount_in" dc:"要兑换的数量，使用最小单位lamports，100000000=0.1SOL" required:"true"`
}

type SolSlippageBpsRes struct {
	g.Meta               `mime:"application/json"`
	ComputedAutoSlippage uint64 `json:"computed_auto_slippage"`
}

// SolSwapRouteReq Solana代币兑换路由请求参数
type SolSwapRouteReq struct {
	g.Meta      `path:"/swapRoute" tags:"区块链" method:"post" summary:"Sol 获取交换路由"`
	TokenIn     string `json:"token_in" dc:"要支付的代币地址 sol地址(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut    string `json:"token_out" dc:"要获得的代币地址" required:"true"`
	AmountIn    uint64 `json:"amount_in" dc:"要兑换的数量，使用最小单位lamports，100000000=0.1SOL" required:"true"`
	SlippageBps uint64 `json:"slippage_bps" dc:"基点滑动容差例如0.01%则填入1，0.03 则填入3" required:"true"`
	PriorityFee uint64 `json:"priority_fee" dc:"优先交易费使用最小单位lamports " binding:"required"`
	SignAddress string `json:"sign_address" required:"true"`
	Adaptor     string `json:"adaptor" dc:"可选值(Jupiter,Raydium)"`
	Mev         bool   `json:"mev" dc:"是否开启防夹"`
}

// SolSwapRouteRes Solana代币兑换路由响应
type SolSwapRouteRes struct {
	g.Meta `mime:"application/json"`
	model.SwapRouterRes
}

type SolSwapQuoteReq struct {
	g.Meta      `path:"/swapRoute" tags:"区块链" method:"post" summary:"Sol 获取交换路由"`
	TokenIn     string `json:"token_in" dc:"要支付的代币地址 sol地址(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut    string `json:"token_out" dc:"要获得的代币地址" required:"true"`
	AmountIn    uint64 `json:"amount_in" dc:"要兑换的数量，使用最小单位lamports，100000000=0.1SOL" required:"true"`
	SlippageBps uint64 `json:"slippage_bps" dc:"基点滑动容差例如0.01%则填入1，0.03 则填入3" required:"true"`
	PriorityFee uint64 `json:"priority_fee" dc:"优先交易费使用最小单位lamports " binding:"required"`
	SignAddress string `json:"sign_address" required:"true"`
	Adaptor     string `json:"adaptor" dc:"可选值(Jupiter,Raydium)"`
	Mev         bool   `json:"mev" dc:"是否开启防夹"`
}

type SolSwapQuoteRes struct {
	g.Meta `mime:"application/json"`
	model.SwapRouterRes
}

type SolSendTxReq struct {
	g.Meta      `path:"/solSendTx" tags:"区块链" method:"post" summary:"Sol 发送交易" deprecated:"true"`
	TxRaw       string `json:"tx_raw" required:"true"`
	SignAddress string `json:"sign_address" required:"true"`
}

type SolSendTxRes struct {
	g.Meta `mime:"application/json"`
	Hex    string `json:"hex" dc:"链上哈希"`
}

type SolGetTxInfoReq struct {
	g.Meta `path:"/solGetTxInfo" tags:"区块链" method:"get" summary:"Sol 获取链上交易信息"`
	Hex    string `json:"hex" dc:"链上哈希"`
}

type SolGetTxInfoRes struct {
	g.Meta `mime:"application/json"`
	*rpc.GetTransactionResult
}

type SolGetMetaDataReq struct {
	g.Meta  `path:"/solGetMetaData" tags:"区块链" method:"get" summary:"Sol 获取token信息"`
	Tokens  string `json:"tokens" dc:"token地址 多个以逗号(,)分隔"`
	Adaptor string `json:"adaptor" dc:"可选值(Jupiter,Raydium)"`
}

type SolGetMetaDataRes struct {
	g.Meta `mime:"application/json"`
	model.TokenInfoRes
}

type SolRpcURLReq struct {
	g.Meta `path:"/solRpcURL" tags:"区块链" method:"get" summary:"Sol RpcURL"`
}

type SolRpcURLRes struct {
	g.Meta `mime:"application/json"`
	URL    []string `json:"url"`
	MevURL string   `json:"mev_url"`
}

type TransactionCallbackReq struct {
	g.Meta      `path:"/transaction-callback" tags:"钱包" method:"post" summary:"交易后回调"`
	Hash        string `json:"hash" dc:"交易hash" required:"true"`
	SignAddress string `json:"sign_address" required:"true"`
}

type TransactionCallbackRes struct {
}
