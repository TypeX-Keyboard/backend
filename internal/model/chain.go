package model

import (
	"keyboard-api-go/internal/util/chain/solana/model"
)

// Quote 交易报价详细信息
type Quote struct {
	InputMint            string        `json:"inputMint" dc:"Enter the Mint address of the token"`
	InAmount             string        `json:"inAmount" dc:"Enter the amount"`
	OutputMint           string        `json:"outputMint" dc:"The Mint address of the output token"`
	OutAmount            string        `json:"outAmount" dc:"Output amount"`
	OtherAmountThreshold string        `json:"otherAmountThreshold" dc:"Other amount thresholds"`
	SwapMode             string        `json:"swapMode" dc:"Exchange mode"`
	SlippageBps          string        `json:"slippageBps" dc:"Slippage base point value"`
	PlatformFee          interface{}   `json:"platformFee" dc:"Platform fees"`
	PriceImpactPct       string        `json:"priceImpactPct" dc:"Price Impact Percentage"`
	RoutePlan            []interface{} `json:"routePlan" dc:"Route plan details"`
	ContextSlot          int           `json:"contextSlot" dc:"Context slots"`
	TimeTaken            float64       `json:"timeTaken" dc:"CONSUMING"`
}

// RawTx 原始交易数据
type RawTx struct {
	SwapTransaction           string `json:"swapTransaction" dc:"Swap serialized data for transactions"`
	LastValidBlockHeight      int    `json:"lastValidBlockHeight" dc:"The height of the last block in which the transaction is valid"`
	PrioritizationFeeLamports int    `json:"prioritizationFeeLamports" dc:"Priority processing fee (in lamports)"`
	RecentBlockhash           string `json:"recentBlockhash" dc:"The most recent block hash"`
}

type SwapRouterRes struct {
	model.GetSwapRouterRes
}

type PriorityFeeRes model.PriorityFeeRes

type TokenInfoRes []model.TokenInfoRes

type TransferRes struct {
	Program interface{} `json:"program"`
	TxRaw   string      `json:"tx_raw"`
}
