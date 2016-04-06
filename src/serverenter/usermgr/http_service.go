/*
* http服务类
 */
package UMGR

import (
	"third/gin"
)

type HttpServer struct {
	*gin.Engine            //http网络通信
}

func CreateHTTPServer() (*HttpServer, error) {
    roter := gin.New()
    return &HttpServer{roter}, nil
}


