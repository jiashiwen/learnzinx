package utils

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"zinx/ziface"
)

/*
全局参数对象，供其他模块使用
*/

type GlobalObj struct {
	TcpServer        ziface.IServer //当前zinx全局Server对象
	Host             string         //当前服务器主机监听的IP
	TcpPort          int            //当前服务器监听的端口
	Name             string         //当前服务器名称
	Version          string
	MaxConn          int    //当前主机允许的最大链接数
	MaxPackageSize   uint32 //Zinx框架数据包最大值
	WorkerPoolSize   uint32 //当前业务worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //zinx框架最多允许的worker数量
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	//配置文件默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServer",
		Version:          "v0.7",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   uint32(runtime.NumCPU()),
		MaxWorkerTaskLen: 1024,
	}

	GlobalObject.Reload()
}
