// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EarningDao is the data access object for table earning.
type EarningDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns EarningColumns // columns contains all the column names of Table for convenient usage.
}

// EarningColumns defines and stores column names for table earning.
type EarningColumns struct {
	Id        string //
	Address   string //
	Usd       string //
	Usd24H    string //
	Cost      string //
	CreatedAt string //
	UpdatedAt string //
}

// earningColumns holds the columns for table earning.
var earningColumns = EarningColumns{
	Id:        "id",
	Address:   "address",
	Usd:       "usd",
	Usd24H:    "usd24h",
	Cost:      "cost",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewEarningDao creates and returns a new DAO object for table data access.
func NewEarningDao() *EarningDao {
	return &EarningDao{
		group:   "default",
		table:   "earning",
		columns: earningColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *EarningDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *EarningDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *EarningDao) Columns() EarningColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *EarningDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *EarningDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *EarningDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
