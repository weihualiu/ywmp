/**
 * 主函数
 */
package main

import (
	"github.com/weihualiu/ywmp/services"
	//_"channelCore"
	_ "channel" //只执行包init，不引入具体包内容
	"time"
	//"test"
)

func main() {
	//启动服务
	services.Start()

	for {
		time.Sleep(100 * time.Second)
	}

	//test.Basic()
}
