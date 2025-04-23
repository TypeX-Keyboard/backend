// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WorkRecordDao is the data access object for table work_record.
type WorkRecordDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns WorkRecordColumns // columns contains all the column names of Table for convenient usage.
}

// WorkRecordColumns defines and stores column names for table work_record.
type WorkRecordColumns struct {
	Id        string //
	Address   string // 公钥地址
	TypeCount string // 本次上报的字符数量
	CreatedAt string //
	UpdatedAt string //
}

// workRecordColumns holds the columns for table work_record.
var workRecordColumns = WorkRecordColumns{
	Id:        "id",
	Address:   "address",
	TypeCount: "type_count",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewWorkRecordDao creates and returns a new DAO object for table data access.
func NewWorkRecordDao() *WorkRecordDao {
	return &WorkRecordDao{
		group:   "default",
		table:   "work_record",
		columns: workRecordColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *WorkRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *WorkRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *WorkRecordDao) Columns() WorkRecordColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *WorkRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *WorkRecordDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *WorkRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
