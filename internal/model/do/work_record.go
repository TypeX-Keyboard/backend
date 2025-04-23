// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WorkRecord is the golang structure of table work_record for DAO operations like Where/Data.
type WorkRecord struct {
	g.Meta    `orm:"table:work_record, do:true"`
	Id        interface{} //
	Address   interface{} // Public key address
	TypeCount interface{} // The number of characters to be reported this time
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
