/*
*配置文件的解析，主要配置存放到全局区
 */
package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var DEFAULT_CONF_FILE string = "./user_manager.conf.default"

type GlobalConfig struct {
	ListenAddr           string
	MysqlConn            string
	MysqlConnectPoolSize int
}

var g_config *GlobalConfig

/*
 *  Description   包的初始化函数，用于创建一个全局的配置文件结构体
 */
func init() {
	g_config = new(GlobalConfig)
}

/*
 *  Description   获取全局配置结构体
 *   Returns         返回获取到的全局配置
 */
func GetGlobalConfig() (*GlobalConfig, error) {
	if g_config == nil {
		return nil, errors.New("全局配置结构体没有被创建即 func init() 调用失败")
	}
	return g_config, nil
}

/*
 *  Description:   初始化全局参数配置
 *  Params       :   全局配置文件的路径
 *   Returns      :   操作成功返回nil, 失败返回具体的error
 */
func (config *GlobalConfig) Init(config_path string) error {
	file, err := os.Open(config_path)
	if err != nil {
		return err
	}
	defer file.Close()

	config_str, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	json.Unmarshal(config_str, config)
	return nil
}
