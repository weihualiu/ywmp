package channel

import "github.com/weihualiu/ywmp/src/services"

func init() {
	//初始化HTTP通道
	userController := UserController{}
	services.SetChannel("/user", &userController)

	ic := indexController{}
	services.SetChannel("/index", &ic)
	//services.SetChannel("/", &ic)

	hbc := heartbeatController{}
	services.SetChannel("/mp/heartbeat", &hbc)

	rg := resourceGenController{}
	services.SetChannel("/resource/", &rg)

}
