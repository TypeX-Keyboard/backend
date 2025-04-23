// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package token

import (
	"context"

	"keyboard-api-go/api/token/v1"
)

type ITokenV1 interface {
	GetTokenDetail(ctx context.Context, req *v1.GetTokenDetailReq) (res *v1.GetTokenDetailRes, err error)
	CheckAddress(ctx context.Context, req *v1.CheckAddressReq) (res *v1.CheckAddressRes, err error)
}
