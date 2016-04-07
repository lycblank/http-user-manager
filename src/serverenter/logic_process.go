/*
* Description 逻辑处理对象，提供服务器的主要逻辑
* 1. 增加用户，监听路径为 POST /user 和 POST /user/:id          可以带参数 id,name, gender, birthday
* 2. 删除用户，监听路径为 DELETE /user 和 DELETE /user/:id   可以带参数 id, name, gender, birthday, low, high
* 3. 更新用户，监听路径为 PUT /user 和 PUT /user/:id               可以带参数 id, name, gender, birthday, low, high
* 4. 查询用户， 监听路径为 GET /user 和 GET /user/:id              可以带参数 id, limit, low, high, name, gender, birthday, offset, order
 */
package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"serverenter/user"
	"strconv"
	"sync/atomic"
	"third/gin"
	"third/gorm"
)

/*
 *  Description:    创建数据库连接池
 *   Returns      :   *DB 连接池指针，　error nil表示成功　非nil表示失败
 */
func CreateDB() (*USER.DB, error) {
	config, err := GetGlobalConfig()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open("mysql", config.MysqlConn)
	if err == nil {
		//同步表结构
		db.AutoMigrate(&USER.User{})
		db.DB().SetMaxOpenConns(config.MysqlConnectPoolSize)
		db.DB().SetMaxIdleConns(config.MysqlConnectPoolSize >> 1)
		return &USER.DB{&db}, nil
	}

	return nil, err
}

type UserManager struct {
	http *HttpServer
	db   USER.DB

	//用于退出服务时，使用的变量
	srv_flag bool  //服务标识，true 表示正常服务， false 表示不进行服务
	srv_num  int32 //标识服务器服务的数量 用于退出服务时使用
}

/*
 *  Description:   初始化用户管理
 *   Returns      :   操作成功返回nil, 失败返回具体的error
 */
func (u_mgr *UserManager) Init() error {
	//1. 初始化全局配置文件
	config, err := GetGlobalConfig()
	if err != nil {
		return err
	}
	err = config.Init(DEFAULT_CONF_FILE)
	if err != nil {
		return err
	}

	//2. 初始化http服务器
	u_mgr.http, err = CreateHTTPServer()
	if err != nil {
		return err
	}

	//3. 初始化db操作
	tmp_db, err := CreateDB()
	if err != nil {
		return err
	}
	u_mgr.db = *tmp_db
	u_mgr.srv_flag = true
	u_mgr.srv_num = 0
	return nil
}

/*
 *  Description:   监听任务操作 ，包括用户的增，删，查，改等操作
 *   Returns      :   操作成功返回nil, 失败返回具体的error
 */
func (u_mgr *UserManager) Start() error {
	//注册增加用户的的操作
	u_mgr.registerAddUserOperation()
	//注册删除用户的操作
	u_mgr.registerDelUserOperation()
	//注册更新用户的操作
	u_mgr.registerUpdateUserOperation()
	//注册查询用户的操作
	u_mgr.registerQueryUserOperation()
	config, err := GetGlobalConfig()
	if err != nil {
		return err
	}
	go u_mgr.http.Run(config.ListenAddr)
	return nil
}

/*
 *  Description:   等待信号，实现服务器的优雅退出
 */
func (u_mgr *UserManager) HandleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	u_mgr.exitFunc()
}

func (u_mgr *UserManager) exitFunc() {
	//服务器标识设置成false 终止服务
	u_mgr.srv_flag = false

	//轮询是否还有服务
	for {
		if atomic.LoadInt32(&u_mgr.srv_num) == 0 {
			//没有服务了 杀死本进程
			os.Exit(0)
		} else {
			//还有服务没有完成， 让出时间片
			runtime.Gosched()
		}
	}
}

func (u_mgr *UserManager) registerUpdateUserOperation() {
	u_mgr.registerUpdateUserByID()
	u_mgr.registerUpdateUser()
}

func (u_mgr *UserManager) updateUser(c *gin.Context) {
	if !u_mgr.canWork() {
		//不能进行工作
		c.JSON(406, gin.H{"status": "服务器关闭中......"})
		return
	}
	//增加服务器的服务数量
	atomic.AddInt32(&u_mgr.srv_num, int32(1))
	defer func() {
		atomic.AddInt32(&u_mgr.srv_num, int32(-1))
	}()

	id_str := c.Param("id")
	var err error
	var usr_pack *USER.UserQueryPack
	if id_str != "" {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "转换id参数错误......"})
			return
		}
		usr_pack, err = u_mgr.getUserPack(c, id)
	} else {
		usr_pack, err = u_mgr.getUserPack(c)
	}

	if err != nil {
		c.JSON(400, gin.H{"error": "获取用户包时,参数错误"})
		return
	}

	if (&usr_pack.Usr).Update(u_mgr.db, usr_pack.IDRange.Low, usr_pack.IDRange.High) != nil {
		c.JSON(405, gin.H{"error": "操作数据库时发生错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"object": usr_pack.Usr})
}

/*
 *  Description:   更新用户通过指定的用户ID, /user/:id
 */
func (u_mgr *UserManager) registerUpdateUserByID() {
	if u_mgr.canWork() {
		u_mgr.http.PUT("/user/:id", func(c *gin.Context) {
			u_mgr.updateUser(c)
		})
	}
}

/*
 *  Description:   注册更新用户操作接口, /user
 */
func (u_mgr *UserManager) registerUpdateUser() {
	if u_mgr.canWork() {
		u_mgr.http.PUT("/user", func(c *gin.Context) {
			u_mgr.updateUser(c)
		})
	}
}

func (u_mgr *UserManager) registerDelUserOperation() {
	u_mgr.registerDelUserByID()
	u_mgr.registerDelUser()
}

func (u_mgr *UserManager) deleteUser(c *gin.Context) {
	if !u_mgr.canWork() {
		//不能进行工作
		c.JSON(406, gin.H{"status": "服务器关闭中......"})
		return
	}

	//增加服务器的服务数量
	atomic.AddInt32(&u_mgr.srv_num, int32(1))
	defer func() {
		atomic.AddInt32(&u_mgr.srv_num, int32(-1))
	}()

	var err error
	var usr_pack *USER.UserQueryPack
	id_str := c.Param("id")
	if id_str != "" {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "转换id参数错误......"})
			return
		}
		usr_pack, err = u_mgr.getUserPack(c, id)
	} else {
		usr_pack, err = u_mgr.getUserPack(c)
	}

	if err != nil {
		c.JSON(400, gin.H{"status": "获取用户包时,参数错误"})
		return
	}
	if (&usr_pack.Usr).Delete(u_mgr.db, usr_pack.IDRange.Low, usr_pack.IDRange.High) != nil {
		c.JSON(405, gin.H{"status": "操作数据库时发生错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"object": usr_pack.Usr})
}

/*
 *  Description:   删除用户通过指定的用户ID, /user/:id
 */
func (u_mgr *UserManager) registerDelUserByID() {
	if u_mgr.canWork() {
		u_mgr.http.DELETE("/user/:id", func(c *gin.Context) {
			u_mgr.deleteUser(c)
		})
	}
}

/*
 *  Description:   注册删除用户操作接口, /user
 */
func (u_mgr *UserManager) registerDelUser() {
	if u_mgr.canWork() {
		u_mgr.http.DELETE("/user", func(c *gin.Context) {
			u_mgr.deleteUser(c)
		})
	}
}

func (u_mgr *UserManager) registerAddUserOperation() {
	u_mgr.registerAddUserByID()
	u_mgr.registerAddUser()
}

func (u_mgr *UserManager) addUser(c *gin.Context) {
	if !u_mgr.canWork() {
		//不能进行工作
		c.JSON(406, gin.H{"status": "服务器关闭中......"})
		return
	}

	//增加服务器的服务数量
	atomic.AddInt32(&u_mgr.srv_num, int32(1))
	defer func() {
		atomic.AddInt32(&u_mgr.srv_num, int32(-1))
	}()

	var err error
	id_str := c.Param("id")
	var usr *USER.User
	if id_str != "" {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "转换id参数错误......"})
			return
		}
		usr, err = u_mgr.getUser(c, id)
	} else {
		usr, err = u_mgr.getUser(c)
	}

	if err != nil {
		c.JSON(400, gin.H{"status": "获取用户包时,参数错误"})
		return
	}
	if usr.Add(u_mgr.db) != nil {
		c.JSON(405, gin.H{"status": "操作数据库时发生错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"object": usr})
}

/*
 *  Description:   增加用户通过指定的用户ID, /user/:id
 */
func (u_mgr *UserManager) registerAddUserByID() {
	if u_mgr.canWork() {
		u_mgr.http.POST("/user/:id", func(c *gin.Context) {
			u_mgr.addUser(c)
		})
	}
}

/*
 *  Description:   注册增加用户操作接口, /user
 */
func (u_mgr *UserManager) registerAddUser() {
	if u_mgr.canWork() {
		u_mgr.http.POST("/user", func(c *gin.Context) {
			u_mgr.addUser(c)
		})
	}
}

func (u_mgr *UserManager) registerQueryUserOperation() {
	u_mgr.registerQueryUserByID()
	u_mgr.registerQueryUser()
}

func (u_mgr *UserManager) queryUser(c *gin.Context) {
	if !u_mgr.canWork() {
		//不能进行工作
		c.JSON(406, gin.H{"status": "服务器关闭中......"})
		return
	}

	//增加服务器的服务数量
	atomic.AddInt32(&u_mgr.srv_num, int32(1))
	defer func() {
		atomic.AddInt32(&u_mgr.srv_num, int32(-1))
	}()

	var err error
	var usr_pack *USER.UserQueryPack
	id_str := c.Param("id")

	if id_str != "" {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "转换id参数错误......"})
			return
		}
		usr_pack, err = u_mgr.getUserPack(c, id)
	} else {
		usr_pack, err = u_mgr.getUserPack(c)
	}

	if err != nil {
		c.JSON(400, gin.H{"status": "获取用户包时,参数错误"})
		return
	}
	usr_list := &USER.UserList{}
	if usr_list.Fetch(u_mgr.db, usr_pack) != nil {
		c.JSON(405, gin.H{"status": "操作数据库时发生错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"object": usr_list})
}

/*
 *  Description:   注册用户通过用户ID进行查询, /user/:id
 */
func (u_mgr *UserManager) registerQueryUserByID() {
	if u_mgr.canWork() {
		u_mgr.http.GET("/user/:id", func(c *gin.Context) {
			u_mgr.queryUser(c)
		})
	}
}

/*
 *  Description:   注册用户查询, /user
 */
func (u_mgr *UserManager) registerQueryUser() {
	if u_mgr.canWork() {
		u_mgr.http.GET("/user", func(c *gin.Context) {
			u_mgr.queryUser(c)
		})
	}
}

/*
 *  Description:   判断该对象是否可以工作
 *  Return        :   true 正常工作， false 不能正常工作
 */
func (u_mgr *UserManager) canWork() bool {
	return u_mgr.http != nil && u_mgr.srv_flag
}

/*
 *  Description:   获取用户包结构体
 *  Param         :  c *gin.Context http服务
 *  Return        :   UserQueryPack 用户结构体  error nil 没有错误， 否则发生错误
 */
func (u_mgr *UserManager) getUserPack(c *gin.Context, ids ...interface{}) (*USER.UserQueryPack, error) {
	usr_pack := USER.UserQueryPack{}
	usr, err := u_mgr.getUser(c, ids...)
	if err != nil {
		return nil, err
	}
	usr_pack.Usr = *usr

	//获取限制
	limit := c.Query("limit")
	usr_pack.Limit = -1
	if limit != "" {
		usr_pack.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}

	//获取排序
	order := c.Query("order")
	usr_pack.Order = 0
	if order != "" {
		usr_pack.Order, err = strconv.Atoi(order)
		if err != nil {
			return nil, err
		}
	}

	//获取偏移
	offset := c.Query("order")
	usr_pack.Offset = -1
	if offset != "" {
		usr_pack.Offset, err = strconv.Atoi(order)
		if err != nil {
			return nil, err
		}
	}

	//获取ID范围
	low := c.Query("low")
	usr_pack.IDRange.Low = -1
	if low != "" {
		usr_pack.IDRange.Low, err = strconv.Atoi(low)
		if err != nil {
			return nil, err
		}
	}

	high := c.Query("high")
	usr_pack.IDRange.High = -1
	if low != "" {
		usr_pack.IDRange.High, err = strconv.Atoi(high)
		if err != nil {
			return nil, err
		}
	}

	return &usr_pack, nil
}

/*
 *  Description:   获取用户结构体
 *  Param         :  c *gin.Context http服务
 *  Return        :   User 用户结构体  error nil 没有错误， 否则发生错误
 */
func (u_mgr *UserManager) getUser(c *gin.Context, ids ...interface{}) (*USER.User, error) {
	usr := USER.User{}
	//获取id
	if len(ids) != 1 {
		var err error
		id := c.Query("id")
		if id != "" {
			usr.ID, err = strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
		}
	} else {
		id, flag := ids[0].(int)
		if !flag {
			return nil, errors.New("断言id值失败，传入的id不是int类型")
		}
		usr.ID = id
	}

	usr.Name = c.Query("name")
	usr.Gender = c.Query("gender")
	usr.Birthday = c.Query("birthday")

	return &usr, nil
}
