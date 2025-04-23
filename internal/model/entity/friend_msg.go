// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FriendMsg is the golang structure for table friend_msg.
type FriendMsg struct {
	Id            int64       `json:"id"            orm:"id"             description:""`
	Address       string      `json:"address"       orm:"address"        description:"Public key address"`
	FriendAddress string      `json:"friendAddress" orm:"friend_address" description:"The friend's public key address"`
	Status        int         `json:"status"        orm:"status"         description:"3=Rejected, 2=Approved, 1=Application"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`
}
