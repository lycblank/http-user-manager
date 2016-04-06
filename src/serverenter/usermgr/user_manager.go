/*
 用户数据库管理模板
 */
package UMGR

import (
        "third/gorm"
        _ "third/go-sql-driver/mysql"
      "errors"
        
)

//用户结构体，存入数据库中的结构
type User struct {
	ID		int              `gorm:"primary_key"`
	Name    	               string
	Gender          	string
	Birthday 	string
}

//范围 [low, high]
type Range struct {
       Low     int
       High    int
}

//查询条件
type UserQueryPack struct {
    Usr         User
    Offset     int           //偏移多少条
    Limit       int          //限制返回多少条记录
    Order     int          // -1 降序 1 升序 
    IDRange Range    //查询ID的范围 
}

/*
 *  Description:    初始化数据库中的表名  
 *   Returns      :   返回数据库中的表名字符串
 */
func (u User) TableName() string {
    return "user"
}

//用户管理类，实现了增加，删除，修改，查询等操作
type DB struct {
    *gorm.DB
}

/*
 *  Description:    创建数据库连接池 
 *   Returns      :   *DB 连接池指针，　error nil表示成功　非nil表示失败
 */
func CreateDB() (*DB, error) {
    config, err := GetGlobalConfig()
    if err != nil {
        return nil, err
    }    
    db, err := gorm.Open("mysql", config.MysqlConn)
    if err == nil {
        //同步表结构
        db.AutoMigrate(&User{})
        db.DB().SetMaxOpenConns(config.MysqlConnectPoolSize)
        db.DB().SetMaxIdleConns(config.MysqlConnectPoolSize >> 1)
        return &DB{&db}, nil
    }
    
    return nil, err
}


func (db *DB) Search(u_pack *UserQueryPack) (usr_slice []User, err error) {
    if u_pack == nil {
        err = errors.New("用户查询包为空 u_pack == nil")
        return
    }
    
    usr := u_pack.Usr
    
    var db_tmp *gorm.DB = db.New()
    
    //数据ID范围限制
    if u_pack.IDRange.Low != -1 && u_pack.IDRange.High != -1 {
        //ID范围有效 u_pack.Usr.ID 强制赋值成0
        usr.ID = 0
        db_tmp = db_tmp.Where("id >= ? and id <= ?", u_pack.IDRange.Low, u_pack.IDRange.High)
    }
    
    //数据偏移
    if u_pack.Offset != -1 {
        db_tmp = db_tmp.Offset(u_pack.Offset)
    }
    
    //数量限制
    if u_pack.Limit != -1 {
        db_tmp = db_tmp.Limit(u_pack.Limit)
    }
    
    //数据排序
    if u_pack.Order != 0 {
         db_tmp = db_tmp.Order("id")
    }
    
    err = db_tmp.Where(&usr).Find(&usr_slice).Error
    return
}

func (db *DB) Del(u_pack *UserQueryPack) error {
      if u_pack == nil {
        return errors.New("用户查询包为空 u_pack == nil")
    }
    
    usr := u_pack.Usr
    
    var db_tmp *gorm.DB = db.New()
    
    //数据ID范围限制
    if u_pack.IDRange.Low != -1 && u_pack.IDRange.High != -1 {
        //ID范围有效 u_pack.Usr.ID 强制赋值成0
        usr.ID = 0
        db_tmp = db_tmp.Where("id >= ? and id <= ?", u_pack.IDRange.Low, u_pack.IDRange.High)
    }
    return db_tmp.Where(&usr).Delete(&usr).Error
}


func (db *DB) Modify(u_pack *UserQueryPack) error {
      if u_pack == nil {
        return errors.New("用户查询包为空 u_pack == nil")
    }
    
    usr := u_pack.Usr
    
    var db_tmp *gorm.DB = db.New().Model(&User{})
    
    //数据ID范围限制
    if u_pack.IDRange.Low != -1 && u_pack.IDRange.High != -1 {
        //ID范围有效 u_pack.Usr.ID 强制赋值成0
        usr.ID = 0
        db_tmp = db_tmp.Where("id >= ? and id <= ?", u_pack.IDRange.Low, u_pack.IDRange.High)
    }
    
    err := db_tmp.Updates(&usr).Error
    return err
}









