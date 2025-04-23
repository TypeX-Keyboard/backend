// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Earning is the golang structure for table earning.
type Earning struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	Address   string      `json:"address"   orm:"address"    description:""`
	Usd       float64     `json:"usd"       orm:"usd"        description:""`
	Usd24H    float64     `json:"usd24H"    orm:"usd24h"     description:""`
	Cost      float64     `json:"cost"      orm:"cost"       description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
