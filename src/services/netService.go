
package services

/*
 * 处理不同访问路径，并建立不同的端口监听
 * 支持多端口多路径控制
 * 端口SSL、Body独立加密TLS
 */

import "log"
import "net/http"
//import "fmt"

/*
 * 功能接口
 */
type ChannelInf interface {
    Process(w http.ResponseWriter, req *http.Request)  //功能核心函数
}

type Channel map[string]ChannelInf

type NetService struct {
    listen uint32  //监听端口
    staticPath string //静态资源路径
    channels Channel
}


/*
 * 启动网络服务
 */
func (ns *NetService) Start() bool {
    
    go func(){
        //检查是否端口是否被占用
        //关联通道
        var cha ChannelInf
        //循环加载通道
        for k,v := range ns.channels {
            cha = v
            http.HandleFunc(k, cha.Process)
        }
        //fmt.Println(cha)
        //启动端口
        log.Fatal(http.ListenAndServe(":9827", nil))
    }()
    return true
}

/*
 * 关闭网络服务
 */
func (ns *NetService) Stop() bool {
    return true
}

/*
 * 设置通道
 */
func (ns *NetService) SetChannel(name string, channel ChannelInf) bool {
    if ns.channels == nil {
        ns.channels = make(Channel, 100)
    }
    ns.channels[name] = channel
    return true
}
