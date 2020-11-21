package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
全局参数对象，供其他模块使用
*/

type GlobalObj struct {
	TcpServer      ziface.IServer //当前zinx全局Server对象
	Host           string         //当前服务器主机监听的IP
	TcpPort        int            //当前服务器监听的端口
	Name           string         //当前服务器名称
	Version        string
	MaxConn        int    //当前主机允许的最大链接数
	MaxPackageSize uint32 //Zinx框架数据包最大值
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
		Name:           "ZinxServer",
		Version:        "v0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalObject.Reload()
}
