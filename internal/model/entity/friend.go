// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Friend is the golang structure for table friend.
type Friend struct {
	Id            int64       `json:"id"            orm:"id"             description:""`
	Address       string      `json:"address"       orm:"address"        description:"Public key address"`
	FriendAddress string      `json:"friendAddress" orm:"friend_address" description:"The friend's public key address"`
	Active        int         `json:"active"        orm:"active"         description:"1 = active, 0 = inactive"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`
}
