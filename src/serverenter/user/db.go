package USER

import (
	_ "third/go-sql-driver/mysql"
	"third/gorm"
)

//用户管理类，实现了增加，删除，修改，查询等操作
type DB struct {
	*gorm.DB
}
