// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package wallet

import (
	"context"

	"keyboard-api-go/api/wallet/v1"
)

type IWalletV1 interface {
	GetSupportedChains(ctx context.Context, req *v1.GetSupportedChainsReq) (res *v1.GetSupportedChainsRes, err error)
	GetTransactionsByAddress(ctx context.Context, req *v1.GetTransactionsByAddressReq) (res *v1.GetTransactionsByAddressRes, err error)
	GetTransactionByHash(ctx context.Context, req *v1.GetTransactionByHashReq) (res *v1.GetTransactionByHashRes, err error)
	TotalValueByAddress(ctx context.Context, req *v1.TotalValueByAddressReq) (res *v1.TotalValueByAddressRes, err error)
	AllTokenBalancesByAddress(ctx context.Context, req *v1.AllTokenBalancesByAddressReq) (res *v1.AllTokenBalancesByAddressRes, err error)
	TokenBalancesByAddress(ctx context.Context, req *v1.TokenBalancesByAddressReq) (res *v1.TokenBalancesByAddressRes, err error)
}
