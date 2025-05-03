// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/util/chain/okx"
)

type (
	IWallet interface {
		// 记录每日以实现收益
		Realized(ctx context.Context)
		SyncHold(ctx context.Context, address string, byAddress []okx.TokenAsset) (map[string]model.TokenObj, error)
		BuildTokenAsset(ctx context.Context, address string, balances []okx.TokenAsset) []model.TokenAsset
		Earnings(ctx context.Context, address string, byAddress *okx.AllTokenBalancesByAddressRes) (earnings float64, earningsRate int, DailyEarnings float64, dailyEarningsRate int, err error)
		BuildHashHistory(ctx context.Context, tx model.TransactionsByAddress, walletAddress string) model.TransactionsByAddress
	}
)

var (
	localWallet IWallet
)

func Wallet() IWallet {
	if localWallet == nil {
		panic("implement not found for interface IWallet, forgot register?")
	}
	return localWallet
}

func RegisterWallet(i IWallet) {
	localWallet = i
}
