// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemSettingDao is the data access object for table system_setting.
type SystemSettingDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns SystemSettingColumns // columns contains all the column names of Table for convenient usage.
}

// SystemSettingColumns defines and stores column names for table system_setting.
type SystemSettingColumns struct {
	Id         string //
	SettingKey string //
	Value      string //
	CreatedAt  string //
	UpdatedAt  string //
}

// systemSettingColumns holds the columns for table system_setting.
var systemSettingColumns = SystemSettingColumns{
	Id:         "id",
	SettingKey: "setting_key",
	Value:      "value",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewSystemSettingDao creates and returns a new DAO object for table data access.
func NewSystemSettingDao() *SystemSettingDao {
	return &SystemSettingDao{
		group:   "default",
		table:   "system_setting",
		columns: systemSettingColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SystemSettingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SystemSettingDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SystemSettingDao) Columns() SystemSettingColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SystemSettingDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SystemSettingDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SystemSettingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
