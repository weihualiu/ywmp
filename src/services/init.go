/*
 * 初始化一些服务的基本环境数据
 */

package services

import "github.com/weihualiu/ywmp/src/conf"

var NetSrv *NetService

// configure file
var Config *conf.Conf

func init() {
	NetSrv = new(NetService)
	conf, err := configSet()
	if err != nil {
		panic(err)
	}
	Config = conf
}

func SetChannel(name string, channel ChannelInf) bool {
	//fmt.Println("init setChannel function")
	return NetSrv.SetChannel(name, channel)
}

func Start() bool {
	return NetSrv.Start()
}

func configSet() (*conf.Conf, error) {
	loader := conf.NewLoader()
	loader.Env().Argv()
	loader.File("conf/config.json")
	c, err := loader.Load()
	if err != nil {
		return nil, err
	}
	return c, nil
}
