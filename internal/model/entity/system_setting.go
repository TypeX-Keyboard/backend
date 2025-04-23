// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemSetting is the golang structure for table system_setting.
type SystemSetting struct {
	Id         int64       `json:"id"         orm:"id"          description:""`
	SettingKey string      `json:"settingKey" orm:"setting_key" description:""`
	Value      string      `json:"value"      orm:"value"       description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
