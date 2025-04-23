// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package chain

import (
	"context"

	"keyboard-api-go/api/chain/v1"
)

type IChainV1 interface {
	SolTransfer(ctx context.Context, req *v1.SolTransferReq) (res *v1.SolTransferRes, err error)
	SolPriorityFee(ctx context.Context, req *v1.SolPriorityFeeReq) (res *v1.SolPriorityFeeRes, err error)
	SolSlippageBps(ctx context.Context, req *v1.SolSlippageBpsReq) (res *v1.SolSlippageBpsRes, err error)
	SolSwapRoute(ctx context.Context, req *v1.SolSwapRouteReq) (res *v1.SolSwapRouteRes, err error)
	SolSwapQuote(ctx context.Context, req *v1.SolSwapQuoteReq) (res *v1.SolSwapQuoteRes, err error)
	SolSendTx(ctx context.Context, req *v1.SolSendTxReq) (res *v1.SolSendTxRes, err error)
	SolGetTxInfo(ctx context.Context, req *v1.SolGetTxInfoReq) (res *v1.SolGetTxInfoRes, err error)
	SolGetMetaData(ctx context.Context, req *v1.SolGetMetaDataReq) (res *v1.SolGetMetaDataRes, err error)
	SolRpcURL(ctx context.Context, req *v1.SolRpcURLReq) (res *v1.SolRpcURLRes, err error)
	TransactionCallback(ctx context.Context, req *v1.TransactionCallbackReq) (res *v1.TransactionCallbackRes, err error)
}
