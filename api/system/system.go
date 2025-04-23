// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package system

import (
	"context"

	"keyboard-api-go/api/system/v1"
)

type ISystemV1 interface {
	Settings(ctx context.Context, req *v1.SettingsReq) (res *v1.SettingsRes, err error)
}
