// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FriendDao is the data access object for table friend.
type FriendDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns FriendColumns // columns contains all the column names of Table for convenient usage.
}

// FriendColumns defines and stores column names for table friend.
type FriendColumns struct {
	Id            string //
	Address       string // 公钥地址
	FriendAddress string // 好友公钥地址
	Active        string // 1=活跃，0=不活跃
	CreatedAt     string //
	UpdatedAt     string //
}

// friendColumns holds the columns for table friend.
var friendColumns = FriendColumns{
	Id:            "id",
	Address:       "address",
	FriendAddress: "friend_address",
	Active:        "active",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewFriendDao creates and returns a new DAO object for table data access.
func NewFriendDao() *FriendDao {
	return &FriendDao{
		group:   "default",
		table:   "friend",
		columns: friendColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FriendDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FriendDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FriendDao) Columns() FriendColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *FriendDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FriendDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FriendDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
