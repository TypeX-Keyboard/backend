// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"

	"github.com/gagliardetto/solana-go/rpc"
)

type (
	IChain interface {
		SyncSolPrice(ctx context.Context)
		// TxCallback 函数用于处理交易回调
		TxCallback(ctx context.Context, hash string, SignAddress string) error
		TxListen(ctx context.Context)
		SyncTx(ctx context.Context)
		TxSuccessListen(ctx context.Context)
		SyncTxHandle(ctx context.Context)
		Hold(ctx context.Context, txRecord *entity.TransactionRecord) error
		// WeightedAverage Calculate the weighted average
		WeightedAverage(values []float64, weights []float64) (float64, error)
		SolTransfer(ctx context.Context, from string, to string, token string, amount uint64) (res *model.TransferRes, err error)
		SolPriorityFee(ctx context.Context) (model.PriorityFeeRes, error)
		SolSlippageBps(ctx context.Context, tokenIn string, tokenOut string, amount uint64) (uint64, error)
		SolSwapRoute(ctx context.Context, adaptor string, signAddress string, tokenIn string, tokenOut string, amountIn uint64, slippageBps uint64, priorityFee uint64, mev bool) (*model.SwapRouterRes, error)
		SolGetMetaData(ctx context.Context, adaptor string, tokens string) (res model.TokenInfoRes, err error)
		SolSendTx(ctx context.Context, txRaw string, signAddress string) (string, error)
		SolGetTxInfo(ctx context.Context, hex string) (out *rpc.GetTransactionResult, err error)
		SolRpcURL(ctx context.Context) []string
	}
)

var (
	localChain IChain
)

func Chain() IChain {
	if localChain == nil {
		panic("implement not found for interface IChain, forgot register?")
	}
	return localChain
}

func RegisterChain(i IChain) {
	localChain = i
}
