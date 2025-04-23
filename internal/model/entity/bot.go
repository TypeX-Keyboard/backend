// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Bot is the golang structure for table bot.
type Bot struct {
	Id         int64       `json:"id"         orm:"id"          description:""`
	Address    string      `json:"address"    orm:"address"     description:"publicKeyAddress"`
	DeviceId   string      `json:"deviceId"   orm:"device_id"   description:""`
	Amount     int64       `json:"amount"     orm:"amount"      description:"Number of EVEs"`
	TypeRate   int         `json:"typeRate"   orm:"type_rate"   description:"Typing progress"`
	TypeCount  int64       `json:"typeCount"  orm:"type_count"  description:"Total number of typing"`
	Active     int         `json:"active"     orm:"active"      description:"1 = active, 0 = inactive"`
	Power      int         `json:"power"      orm:"power"       description:"Current power"`
	PowerLimit int         `json:"powerLimit" orm:"power_limit" description:"Power cap"`
	AutoAccept int         `json:"autoAccept" orm:"auto_accept" description:"1=Automatically accepted"`
	Buff       int         `json:"buff"       orm:"buff"        description:"Gain (in %)"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
