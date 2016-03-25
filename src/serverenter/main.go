
package main

import (
	"fmt"
	"serverenter/usermgr"
	"serverenter/net"
	"github.com/gin-gonic/gin"
	//"log"
	//"fmt"
)

//声明回调函数 具体的回调函数声明，参考httpservice.go文件
func Query(usrs []UserManager.User, c *gin.Context, status int) {
	if status != -1 {
		for _, usr := range usrs {
			fmt.Println(usr)
		}
	}
}

func main() {
	httpnet := new(UserHttpServer.UserHttpServer)
	//1. 初始化
	if httpnet.Init() != 0 {
		return
	}
	//2. 注册路径
	httpnet.RegisterQuery("/query", nil, Query)
	httpnet.RegisterModify("/update", nil, nil)
	httpnet.RegisterDel("/delete", nil, nil)
	httpnet.RegisterAdd("/add", nil, nil)
	//3. 监听
	httpnet.Run("127.0.0.1", 8080)
}


