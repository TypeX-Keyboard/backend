// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Earning is the golang structure of table earning for DAO operations like Where/Data.
type Earning struct {
	g.Meta    `orm:"table:earning, do:true"`
	Id        interface{} //
	Address   interface{} //
	Usd       interface{} //
	Usd24H    interface{} //
	Cost      interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
