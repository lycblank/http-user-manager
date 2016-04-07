/*
 用户数据库管理模板
*/
package USER

import (
	"errors"
)

//用户结构体，存入数据库中的结构
type User struct {
	ID       int `gorm:"primary_key"`
	Name     string
	Gender   string
	Birthday string
}

/*
 *  Description:    初始化数据库中的表名
 *   Returns      :   返回数据库中的表名字符串
 */
func (u User) TableName() string {
	return "user"
}

func (usr *User) Update(db DB, low, high int) error {

	update := db.Model(&User{})

	//数据ID范围限制
	if low != -1 && high != -1 {
		//ID范围有效 u_pack.Usr.ID 强制赋值成0
		usr.ID = 0
		update = update.Where("id >= ? and id <= ?", low, high)
	}

	err := update.Updates(usr).Error
	return err
}

func (usr *User) Delete(db DB, low, high int) error {

	del := db.Model(&User{})

	//数据ID范围限制
	if low != -1 && high != -1 {
		//ID范围有效 u_pack.Usr.ID 强制赋值成0
		usr.ID = 0
		del = del.Where("id >= ? and id <= ?", low, high)
	}

	return del.Delete(usr).Error
}

func (usr *User) Add(db DB) error {
	add := db.Model(&User{})
	return add.Create(usr).Error
}

//范围 [low, high]
type Range struct {
	Low  int
	High int
}

//查询条件
type UserQueryPack struct {
	Usr     User
	Offset  int   //偏移多少条
	Limit   int   //限制返回多少条记录
	Order   int   // -1 降序 1 升序
	IDRange Range //查询ID的范围
}

type UserList []User

func (usr_list *UserList) Fetch(db DB, u_pack *UserQueryPack) error {
	if u_pack == nil {
		return errors.New("用户查询包为空 u_pack == nil")
	}

	usr := u_pack.Usr

	query := db.Model(&User{})

	//数据ID范围限制
	if u_pack.IDRange.Low != -1 && u_pack.IDRange.High != -1 {
		//ID范围有效 u_pack.Usr.ID 强制赋值成0
		usr.ID = 0
		query = query.Where("id >= ? and id <= ?", u_pack.IDRange.Low, u_pack.IDRange.High)
	}

	//数据偏移
	if u_pack.Offset != -1 {
		query = query.Offset(u_pack.Offset)
	}

	//数量限制
	if u_pack.Limit != -1 {
		query = query.Limit(u_pack.Limit)
	}

	//数据排序
	if u_pack.Order != 0 {
		query = query.Order("id")
	}

	return query.Where(&usr).Find(usr_list).Error
}
