/*
* 程序启动入口
*/
package main

import (
	"serverenter/usermgr"
)

func main() {
    usr_manager := new(UMGR.UserManager)
    usr_manager.Init()
    usr_manager.Start()
    usr_manager.HandleSignals()
}


