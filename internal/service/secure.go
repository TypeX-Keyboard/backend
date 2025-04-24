// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ISecure interface {
		AcquirePublicKey(ctx context.Context) (string, error)
		SubmitKey(ctx context.Context, uuid string, key string) error
		GenClientKey(ctx context.Context) error
		InitClientKey(ctx context.Context) error
	}
)

var (
	localSecure ISecure
)

func Secure() ISecure {
	if localSecure == nil {
		panic("implement not found for interface ISecure, forgot register?")
	}
	return localSecure
}

func RegisterSecure(i ISecure) {
	localSecure = i
}
