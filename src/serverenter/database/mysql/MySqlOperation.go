/*
 数据库具体操作类(MySql)
 */

package MySqlOperation

import (
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mssql"
	_"github.com/jinzhu/gorm/dialects/mysql"
	_"github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	//"strings"
	"database/sql"
)

//Mysql数据库管理
type MySqlOperation struct {
	DBName 		string
	TableName 	string
	DB			*gorm.DB
}


//实现DatabaseManager接口
func (this *MySqlOperation) Open(username, password, dbname string) error {
	var err error
	this.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbname))
	if err != nil {
		return err
	}
	//默认设置最大空闲数为10
	this.DB.DB().SetMaxIdleConns(10)

   return err
}

//关闭数据库
func (this *MySqlOperation) Close() error {
	return this.DB.Close()
}

//设置最大的空闲连接数
func (this *MySqlOperation) SetMaxIdleConns(idleConn int) error {
	this.DB.DB().SetMaxIdleConns(idleConn)
	return nil
}

//设置最大的打开连接的数
func (this *MySqlOperation) SetMaxOpenConns(openConn int) error {
	this.DB.DB().SetMaxOpenConns(openConn)
	return nil
}

//增加列
func (this *MySqlOperation) AddColumn(col string) error {
	return nil
}

//删除列
func (this *MySqlOperation) DelColumn(col string) error {
	return nil
}

//修改列
func (this *MySqlOperation) ModifyColumn(oldCol, newCol string) error {
	return nil
}

//增加记录
func (this *MySqlOperation) AddRecord(record interface{}) error {
	//构造查询字符串
	return this.DB.Create(record).Error
}

//更新记录
func (this *MySqlOperation) UpdateRecord(record interface{}) error {
	return this.DB.Save(record).Error
}

//查询记录
func (this *MySqlOperation) QueryRecord(where interface{}, args ...interface{})(rows *sql.Rows, err error) {
	return this.DB.Table(this.TableName).Where(where, args...).Rows()	
}

//删除记录
func (this *MySqlOperation) DelRecord(value interface{}) error {
	return this.DB.Delete(value).Error
}

//获取原始的db
func (this *MySqlOperation) GetDB() (db *gorm.DB, err error){
	return this.DB, this.DB.Error
}

//返回数据库中的错误信息
func (this *MySqlOperation) Error() error {
	return this.DB.Error
}





