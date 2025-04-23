// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FriendMsgDao is the data access object for table friend_msg.
type FriendMsgDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns FriendMsgColumns // columns contains all the column names of Table for convenient usage.
}

// FriendMsgColumns defines and stores column names for table friend_msg.
type FriendMsgColumns struct {
	Id            string //
	Address       string // 公钥地址
	FriendAddress string // 好友公钥地址
	Status        string // 3=已拒绝、2=已通过、1=申请中
	CreatedAt     string //
	UpdatedAt     string //
}

// friendMsgColumns holds the columns for table friend_msg.
var friendMsgColumns = FriendMsgColumns{
	Id:            "id",
	Address:       "address",
	FriendAddress: "friend_address",
	Status:        "status",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewFriendMsgDao creates and returns a new DAO object for table data access.
func NewFriendMsgDao() *FriendMsgDao {
	return &FriendMsgDao{
		group:   "default",
		table:   "friend_msg",
		columns: friendMsgColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FriendMsgDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FriendMsgDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FriendMsgDao) Columns() FriendMsgColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *FriendMsgDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FriendMsgDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FriendMsgDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
