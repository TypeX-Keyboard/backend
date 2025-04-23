// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Hold is the golang structure of table hold for DAO operations like Where/Data.
type Hold struct {
	g.Meta       `orm:"table:hold, do:true"`
	Id           interface{} //
	Address      interface{} // walletAddress
	TokenAddress interface{} // tokenAddress
	Amount       interface{} // Number of positions
	CostPrice    interface{} // Cost price
	Cost         interface{} // Settled costs
	Earning      interface{} // Realized gains
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
