// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Bot is the golang structure of table bot for DAO operations like Where/Data.
type Bot struct {
	g.Meta     `orm:"table:bot, do:true"`
	Id         interface{} //
	Address    interface{} // publicKeyAddress
	DeviceId   interface{} //
	Amount     interface{} // number of eves
	TypeRate   interface{} // typing progress
	TypeCount  interface{} // total number of typing
	Active     interface{} // 1 = active, 0 = inactive
	Power      interface{} // Current power
	PowerLimit interface{} // Power cap
	AutoAccept interface{} // 1=Auto-accept
	Buff       interface{} // Gain (in %)
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
