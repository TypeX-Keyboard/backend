package v1

import (
	"keyboard-api-go/internal/model"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gogf/gf/v2/frame/g"
)

type SolTransferReq struct {
	g.Meta `path:"/transfer" tags:"Blockchain" method:"post" summary:"Sol Tx"`
	From   string `json:"from" dc:"Sending address" required:"true"`
	To     string `json:"to" dc:"Receiving address" required:"true"`
	Token  string `json:"token" dc:"Token address sol(So11111111111111111111111111111111111111112)" required:"true"`
	Amount uint64 `json:"amount" dc:"The amount of transfer needs to be multiplied by the decimal of tokens" required:"true"`
}

type SolTransferRes struct {
	g.Meta `mime:"application/json"`
	model.TransferRes
}

type SolPriorityFeeReq struct {
	g.Meta `path:"/solPriorityFee" tags:"Blockchain" method:"get" summary:"Sol Get priority transaction fees"`
}

type SolPriorityFeeRes struct {
	g.Meta `mime:"application/json"`
	model.PriorityFeeRes
}

type SolSlippageBpsReq struct {
	g.Meta   `path:"/solSlippageBps" tags:"Blockchain" method:"get" summary:"Sol Get dynamic slippage"`
	TokenIn  string `json:"token_in" dc:"The address of the token to be paid sol(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut string `json:"token_out" dc:"The address of the token to be obtained" required:"true"`
	AmountIn uint64 `json:"amount_in" dc:"The quantity to be redeemed, using the smallest unit lamports，100000000=0.1SOL" required:"true"`
}

type SolSlippageBpsRes struct {
	g.Meta               `mime:"application/json"`
	ComputedAutoSlippage uint64 `json:"computed_auto_slippage"`
}

type SolSwapRouteReq struct {
	g.Meta      `path:"/swapRoute" tags:"Blockchain" method:"post" summary:"Sol Get switched routes"`
	TokenIn     string `json:"token_in" dc:"The address of the token to be paid sol(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut    string `json:"token_out" dc:"The address of the token to be obtained" required:"true"`
	AmountIn    uint64 `json:"amount_in" dc:"The quantity to be redeemed, using the smallest unit lamports，100000000=0.1SOL" required:"true"`
	SlippageBps uint64 `json:"slippage_bps" dc:"For example, 0.01% is 1 and 0.03 is 3" required:"true"`
	PriorityFee uint64 `json:"priority_fee" dc:"Priority transaction fees use the smallest units lamports " binding:"required"`
	SignAddress string `json:"sign_address" required:"true"`
	Adaptor     string `json:"adaptor" dc:"Optional(Jupiter,Raydium)"`
	Mev         bool   `json:"mev" dc:"Whether to turn on the anti-pinch guard"`
}

type SolSwapRouteRes struct {
	g.Meta `mime:"application/json"`
	model.SwapRouterRes
}

type SolSwapQuoteReq struct {
	g.Meta      `path:"/swapRoute" tags:"Blockchain" method:"post" summary:"Sol Get switched routes"`
	TokenIn     string `json:"token_in" dc:"The address of the token to be paid sol(So11111111111111111111111111111111111111112)" required:"true"`
	TokenOut    string `json:"token_out" dc:"The address of the token to be obtained" required:"true"`
	AmountIn    uint64 `json:"amount_in" dc:"The amount to be exchanged, using the smallest unit of lamports, 100000000 = 0.1 SOL" required:"true"`
	SlippageBps uint64 `json:"slippage_bps" dc:"For example, 0.01% is 1 and 0.03 is 3" required:"true"`
	PriorityFee uint64 `json:"priority_fee" dc:"The priority transaction fee uses the smallest unit, lamports " binding:"required"`
	SignAddress string `json:"sign_address" required:"true"`
	Adaptor     string `json:"adaptor" dc:"Optional(Jupiter,Raydium)"`
	Mev         bool   `json:"mev" dc:"Whether to turn on the anti-pinch guard"`
}

type SolSwapQuoteRes struct {
	g.Meta `mime:"application/json"`
	model.SwapRouterRes
}

type SolSendTxReq struct {
	g.Meta      `path:"/solSendTx" tags:"Blockchain" method:"post" summary:"Sol Send a transaction" deprecated:"true"`
	TxRaw       string `json:"tx_raw" required:"true"`
	SignAddress string `json:"sign_address" required:"true"`
}

type SolSendTxRes struct {
	g.Meta `mime:"application/json"`
	Hex    string `json:"hex" dc:"On-chain hashing"`
}

type SolGetTxInfoReq struct {
	g.Meta `path:"/solGetTxInfo" tags:"Blockchain" method:"get" summary:"Sol Get on-chain transaction information"`
	Hex    string `json:"hex" dc:"On-chain hashing"`
}

type SolGetTxInfoRes struct {
	g.Meta `mime:"application/json"`
	*rpc.GetTransactionResult
}

type SolGetMetaDataReq struct {
	g.Meta  `path:"/solGetMetaData" tags:"Blockchain" method:"get" summary:"Sol Obtain the token information"`
	Tokens  string `json:"tokens" dc:"Token addresses are separated by a comma (,)."`
	Adaptor string `json:"adaptor" dc:"Optional(Jupiter,Raydium)"`
}

type SolGetMetaDataRes struct {
	g.Meta `mime:"application/json"`
	model.TokenInfoRes
}

type SolRpcURLReq struct {
	g.Meta `path:"/solRpcURL" tags:"Blockchain" method:"get" summary:"Sol RpcURL"`
}

type SolRpcURLRes struct {
	g.Meta `mime:"application/json"`
	URL    []string `json:"url"`
	MevURL string   `json:"mev_url"`
}

type TransactionCallbackReq struct {
	g.Meta      `path:"/transaction-callback" tags:"wallet" method:"post" summary:"post trade callbacks"`
	Hash        string `json:"hash" dc:"transaction hash" required:"true"`
	SignAddress string `json:"sign_address" required:"true"`
}

type TransactionCallbackRes struct {
}
