/**
 * 主函数
 */
package main

import (
	"github.com/weihualiu/ywmp/src/services"
	//_"channelCore"
	_ "github.com/weihualiu/ywmp/src/channel" //只执行包init，不引入具体包内容
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
