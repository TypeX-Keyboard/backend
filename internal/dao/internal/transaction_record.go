// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TransactionRecordDao is the data access object for table transaction_record.
type TransactionRecordDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns TransactionRecordColumns // columns contains all the column names of Table for convenient usage.
}

// TransactionRecordColumns defines and stores column names for table transaction_record.
type TransactionRecordColumns struct {
	Id          string //
	Hash        string // 交易哈希
	TokenIn     string // 要支付的代币地址
	TokenOut    string // 要获得的代币地址
	AmountIn    string // 要兑换的数量
	AmountOut   string // 获得代币的数量
	SolPrice    string //
	SignAddress string // 签名的钱包地址
	Status      string // 0=处理中，1=成功，2=失败
	Action      string // 1=兑换，2=转账
	IsHandle    string //
	UpdatedAt   string //
	CreatedAt   string //
}

// transactionRecordColumns holds the columns for table transaction_record.
var transactionRecordColumns = TransactionRecordColumns{
	Id:          "id",
	Hash:        "hash",
	TokenIn:     "token_in",
	TokenOut:    "token_out",
	AmountIn:    "amount_in",
	AmountOut:   "amount_out",
	SolPrice:    "sol_price",
	SignAddress: "sign_address",
	Status:      "status",
	Action:      "action",
	IsHandle:    "is_handle",
	UpdatedAt:   "updated_at",
	CreatedAt:   "created_at",
}

// NewTransactionRecordDao creates and returns a new DAO object for table data access.
func NewTransactionRecordDao() *TransactionRecordDao {
	return &TransactionRecordDao{
		group:   "default",
		table:   "transaction_record",
		columns: transactionRecordColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TransactionRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TransactionRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TransactionRecordDao) Columns() TransactionRecordColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TransactionRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TransactionRecordDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TransactionRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
