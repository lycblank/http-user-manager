/*
 提供数据库接口, 主要用于用户管理时使用
 */


package UserDatabase

import (
	"github.com/jinzhu/gorm"
	"database/sql"
)


//数据库操作接口
type UserDatabase interface {
	//打开数据库
	Open(username, password, dbname string) error
	//关闭数据库
	Close() error
	//增加记录
	AddRecord(record interface{}) error
	//更新记录
	UpdateRecord(record interface{}) error
	//查询记录
	QueryRecord(where interface{}, args ...interface{})(rows *sql.Rows, err error)
	//删除记录
	DelRecord(value interface{}) error
	//获取原始数据库管理对象
	GetDB() (db *gorm.DB, err error)
	//获取错误码
	Error() error
}


