// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TransactionRecord is the golang structure for table transaction_record.
type TransactionRecord struct {
	Id          int64       `json:"id"          orm:"id"           description:""`
	Hash        string      `json:"hash"        orm:"hash"         description:"Transaction hash"`
	TokenIn     string      `json:"tokenIn"     orm:"token_in"     description:"The address of the token to be paid"`
	TokenOut    string      `json:"tokenOut"    orm:"token_out"    description:"The address of the token to be obtained"`
	AmountIn    float64     `json:"amountIn"    orm:"amount_in"    description:"The quantity to be redeemed"`
	AmountOut   float64     `json:"amountOut"   orm:"amount_out"   description:"The number of tokens earned"`
	SolPrice    float64     `json:"solPrice"    orm:"sol_price"    description:""`
	SignAddress string      `json:"signAddress" orm:"sign_address" description:"The wallet address of the signature"`
	Status      int         `json:"status"      orm:"status"       description:"0=Processing, 1=Success, 2=Failure"`
	Action      int         `json:"action"      orm:"action"       description:"1 = exchange, 2 = transfer"`
	IsHandle    int         `json:"isHandle"    orm:"is_handle"    description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
}
