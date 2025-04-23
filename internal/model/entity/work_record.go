// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// WorkRecord is the golang structure for table work_record.
type WorkRecord struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	Address   string      `json:"address"   orm:"address"    description:"Public key address"`
	TypeCount int         `json:"typeCount" orm:"type_count" description:"The number of characters to be reported this time"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
