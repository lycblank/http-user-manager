/*
对外主要使用的类，主要提供http（用户的增，删，查，改等服务） 
主要的使用的方法，参考main.go
 */


package UserHttpServer

import (
	"github.com/gin-gonic/gin"
	"serverenter/usermgr"
	"fmt"
	"strconv"
)

//操作进行前回调过滤 true进行过滤 false不进行过滤
type CallbackBefore func(usr *UserManager.User, c *gin.Context) bool
//操作完成后进行回调 status 表示操作状态 0 代表成功
type CallbackAfter func(usr *UserManager.User, c *gin.Context, status int)
type QueryCallbackAfter func(usr []UserManager.User, c *gin.Context, status int)

type UserHttpServer struct {
	router     			*gin.Engine            //http网络通信
	usrmgr          	*UserManager.UserMgr   //用户管理类
	callbackAfterMap     map[string]CallbackAfter //处理完操作后进行回调
	callbackBeforeMap    map[string]CallbackBefore //处理操作之前进行回调
	queryCallbackAfter   QueryCallbackAfter
}

//初始化操作
func (this *UserHttpServer) Init() (status int) {
	this.router = gin.Default()
	this.usrmgr = new(UserManager.UserMgr)
	this.callbackBeforeMap = make(map[string]CallbackBefore, 3)
	this.callbackAfterMap = make(map[string]CallbackAfter, 3)
	if this.usrmgr.Init() != nil {
		return -1
	}
	return 0
}

//设置自定义的http服务引擎
func (this *UserHttpServer) SetEngine(r *gin.Engine) {
	this.router = r
}

//运行http服务器
func (this *UserHttpServer) Run(ip string, port int) {
	if this.router != nil {
		this.router.Run(fmt.Sprintf("%s:%d", ip, port))
	}
}

//注册增加用户时的回调操作
func (this *UserHttpServer) RegisterAdd(path string, callBefore CallbackBefore, callAfter CallbackAfter){
	this.callbackBeforeMap["Add"] = callBefore
	this.callbackAfterMap["Add"] = callAfter
	if this.router != nil {
		this.router.POST(path, func(c *gin.Context) {
			 usr := this.getUser(c)
			 if this.callbackBeforeMap["Add"] != nil && this.callbackBeforeMap["Add"](usr, c) {
				//被应用程序过滤掉 不进行处理
				return;
			}
			go this.add(usr, c)
		})
	}
}

//注册删除用户时的回调操作
func (this *UserHttpServer) RegisterDel(path string, callBefore CallbackBefore, callAfter CallbackAfter) {
	this.callbackBeforeMap["Del"] = callBefore
	this.callbackAfterMap["Del"] = callAfter
	if this.router != nil {
		this.router.DELETE(path, func(c *gin.Context) {
			usr := this.getUser(c)
			if this.callbackBeforeMap["Del"] != nil && this.callbackBeforeMap["Del"](usr, c) {
				//被应用程序过滤掉 不进行处理
				return;
			}
			go this.del(usr, c)
		})
	}
}

//注册修改用户时的回调操作
func (this *UserHttpServer) RegisterModify(path string, callBefore CallbackBefore, callAfter CallbackAfter) {
	this.callbackBeforeMap["Modify"] = callBefore
	this.callbackAfterMap["Modify"] = callAfter
	if this.router != nil {
		this.router.PUT(path, func(c *gin.Context) {
			usr := this.getUser(c)	
			if this.callbackBeforeMap["Modify"] != nil && this.callbackBeforeMap["Modify"](usr, c) {
				//被应用程序过滤掉 不进行处理
					return;
				}
				go this.modify(usr, c)
			})
	}
}

//注册查询用户时的回调操作
func (this *UserHttpServer) RegisterQuery(path string, callBefore CallbackBefore, callAfter QueryCallbackAfter) {
	this.callbackBeforeMap["Query"] = callBefore
	this.queryCallbackAfter = callAfter
	if this.router != nil {
		this.router.GET(path, func(c *gin.Context) {
			usr := this.getUser(c)
			if this.callbackBeforeMap["Query"] != nil && this.callbackBeforeMap["Query"](usr, c) {
				//被应用程序过滤掉 不进行处理
				return;
			}
			go this.query(usr, c)
		})
	}
}

//协程内部增加处理 处理完后，调用用户注册的回调函数
func (this *UserHttpServer) add(usr *UserManager.User, c *gin.Context) {
	if this.usrmgr != nil {
		status := this.usrmgr.Add(usr)
		if this.callbackAfterMap["Add"] != nil {
			this.callbackAfterMap["Add"](usr, c, status)
		}
	}
}

//协程内部删除处理 处理完后，调用用户注册的回调函数
func (this *UserHttpServer) del(usr *UserManager.User, c *gin.Context) {
	if this.usrmgr != nil {
		status := this.usrmgr.Del(usr)
		if this.callbackAfterMap["Del"] != nil {
			this.callbackAfterMap["Del"](usr, c, status)
		}
	}
}

//协程内部修改处理 处理完后，调用用户注册的回调函数
func (this *UserHttpServer) modify(usr *UserManager.User, c *gin.Context) {
	if this.usrmgr != nil {
		status := this.usrmgr.Modify(usr)
		if this.callbackAfterMap["Modify"] != nil {
			this.callbackAfterMap["Modify"](usr, c, status)
		}
	}
}

//协程内部查询处理 处理完后，调用用户注册的回调函数
func (this *UserHttpServer) query(usr *UserManager.User, c *gin.Context) {
	if  this.usrmgr != nil {
		usrvec, status := this.usrmgr.Query(usr)
		if this.queryCallbackAfter != nil {
			this.queryCallbackAfter(usrvec, c, status)
		}
	}
}

//获取http服务中的用户参数
func (this *UserHttpServer) getUser(c *gin.Context) *UserManager.User {
	usr := UserManager.User{}
	usr.ID, _ = strconv.Atoi(c.Query("id"))
	usr.Name = c.Query("name")
	usr.Gender = c.Query("gender")
	usr.Birthday = c.Query("birthday")
	return &usr
}




