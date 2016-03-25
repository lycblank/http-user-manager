/*
 用户管理类, 提供增加删除修改的功能
 */


package UserManager

import (
	"serverenter/database/mysql"
	"serverenter/database/interface"
	"fmt"
)

//用户结构体，存入数据库中的结构
type User struct {
	ID			int   `gorm:"primary_key"`
	Name    	string
	Gender  	string
	Birthday 	string
}

//用户管理类，实现了增加，删除，修改，查询等操作
type UserMgr struct {
	db    UserDatabase.UserDatabase
}

//初始化数据库相关操作
func (u *UserMgr) Init() error {
	//数据库存储默认采用mysql进行存储
	u.SetDatabase(new(MySqlOperation.MySqlOperation))
	err := u.db.Open("test", "123456", "testDB")
	if err != nil {
		return err
	}
	RawDB, _ := u.db.GetDB()
	//关闭表名的复数形式
	RawDB.SingularTable(true)
	//同步表结构
	err = RawDB.AutoMigrate(&User{}).Error
	return err
}

//自定义数据库存储操作
func (u *UserMgr) SetDatabase(db UserDatabase.UserDatabase) {
	u.db = db
}

//添加用户
func (u *UserMgr) Add(usr *User) (status int) {
	if u.db.Error() != nil {
		return -1
	}
	if u.db.AddRecord(usr) != nil {
		return -1
	}
	return 0
}

//删除用户
func (u *UserMgr) Del(usr *User) (status int) {
	if u.db.Error() != nil {
		return -1
	}
	if u.db.DelRecord(usr) != nil {
		return -1
	}
	return 0
}

//修改用户
func (u *UserMgr) Modify(usr *User) (status int) {
	if u.db.Error() != nil {
		return -1
	}
	if u.db.UpdateRecord(usr) != nil {
		return -1
	}
	return 0
}

//查询用户
func (u *UserMgr) Query(query interface{}, args ...interface{}) (usrpkg []User, status int) {
	if err := u.db.Error(); err != nil {
		fmt.Println(err)
		return nil, -1
	}
	RawDB, _ := u.db.GetDB()
	rows, err := RawDB.Table("user").Where(query, args...).Rows()
	if err != nil {
		fmt.Println(err)
		return nil, -1
	}
	//浏览记录
	var usrvec []User
	usr := User{}
	for rows.Next() {
		rows.Scan(&usr.ID, &usr.Name, &usr.Gender, &usr.Birthday)
		usrvec = append(usrvec, usr)
	}
	return  usrvec, 0
}



