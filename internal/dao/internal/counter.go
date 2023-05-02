// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CounterDao is the data access object for table counter.
type CounterDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns CounterColumns // columns contains all the column names of Table for convenient usage.
}

// CounterColumns defines and stores column names for table counter.
type CounterColumns struct {
	Id        string //
	Db        string //
	CreatedAt string //
	Ip        string //
}

// counterColumns holds the columns for table counter.
var counterColumns = CounterColumns{
	Id:        "id",
	Db:        "db",
	CreatedAt: "created_at",
	Ip:        "ip",
}

// NewCounterDao creates and returns a new DAO object for table data access.
func NewCounterDao() *CounterDao {
	return &CounterDao{
		group:   "default",
		table:   "counter",
		columns: counterColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CounterDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CounterDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CounterDao) Columns() CounterColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CounterDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CounterDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CounterDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}