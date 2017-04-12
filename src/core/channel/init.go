package channelCore

import "github/com/weihualiu/ywmp/services"

func init() {
	//初始化HTTP通道
	userHelloController := UserHelloController{}
	services.SetChannel("/user/hello", &userHelloController)
	userExchangeController := UserExchangeController{}
	services.SetChannel("/user/exchange", &userExchangeController)
}
