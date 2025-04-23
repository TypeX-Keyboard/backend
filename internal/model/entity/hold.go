// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Hold is the golang structure for table hold.
type Hold struct {
	Id           int64       `json:"id"           orm:"id"            description:""`
	Address      string      `json:"address"      orm:"address"       description:"Wallet address"`
	TokenAddress string      `json:"tokenAddress" orm:"token_address" description:"Token address"`
	Amount       float64     `json:"amount"       orm:"amount"        description:"Number of positions"`
	CostPrice    float64     `json:"costPrice"    orm:"cost_price"    description:"Cost price"`
	Cost         float64     `json:"cost"         orm:"cost"          description:"Settled costs"`
	Earning      float64     `json:"earning"      orm:"earning"       description:"Realized gains"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
