// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemSetting is the golang structure of table system_setting for DAO operations like Where/Data.
type SystemSetting struct {
	g.Meta     `orm:"table:system_setting, do:true"`
	Id         interface{} //
	SettingKey interface{} //
	Value      interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
