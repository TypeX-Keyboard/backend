// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Friend is the golang structure of table friend for DAO operations like Where/Data.
type Friend struct {
	g.Meta        `orm:"table:friend, do:true"`
	Id            interface{} //
	Address       interface{} // publicKeyAddress
	FriendAddress interface{} // The friend's public key address
	Active        interface{} // 1 = active, 0 = inactive
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
