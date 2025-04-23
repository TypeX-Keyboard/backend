// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TransactionRecord is the golang structure of table transaction_record for DAO operations like Where/Data.
type TransactionRecord struct {
	g.Meta      `orm:"table:transaction_record, do:true"`
	Id          interface{} //
	Hash        interface{} // Transaction hash
	TokenIn     interface{} // The address of the token to be paid
	TokenOut    interface{} // The address of the token to be obtained
	AmountIn    interface{} // The quantity to be redeemed
	AmountOut   interface{} // The number of tokens earned
	SolPrice    interface{} //
	SignAddress interface{} // The wallet address of the signature
	Status      interface{} // 0=Processing, 1=Success, 2=Failure
	Action      interface{} // 1 = exchange, 2 = transfer
	IsHandle    interface{} //
	UpdatedAt   *gtime.Time //
	CreatedAt   *gtime.Time //
}
