/*
* 程序启动入口
 */
package main

func main() {
	usr_manager := new(UserManager)
	usr_manager.Init()
	usr_manager.Start()
	usr_manager.HandleSignals()
}
