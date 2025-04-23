// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FriendMsg is the golang structure of table friend_msg for DAO operations like Where/Data.
type FriendMsg struct {
	g.Meta        `orm:"table:friend_msg, do:true"`
	Id            interface{} //
	Address       interface{} // publicKeyAddress
	FriendAddress interface{} // The friend's public key address
	Status        interface{} // 3=Rejected, 2=Approved, 1=Application
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
