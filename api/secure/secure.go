// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package secure

import (
	"context"

	"keyboard-api-go/api/secure/v1"
)

type ISecureV1 interface {
	Acquire(ctx context.Context, req *v1.AcquireReq) (res *v1.AcquireRes, err error)
	Submit(ctx context.Context, req *v1.SubmitReq) (res *v1.SubmitRes, err error)
	TestGet(ctx context.Context, req *v1.TestGetReq) (res *v1.TestGetRes, err error)
	TestPost(ctx context.Context, req *v1.TestPostReq) (res *v1.TestPostRes, err error)
	GetKey(ctx context.Context, req *v1.GetKeyReq) (res *v1.GetKeyRes, err error)
	GenKey(ctx context.Context, req *v1.GenKeyReq) (res *v1.GenKeyRes, err error)
}
