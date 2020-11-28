package ziface

import "net"

type IConnection interface {
	//启动链接
	Start()

	//停止链接
	Stop()

	//获取当前链接的绑定socket conn
	GetTcpConnection() *net.TCPConn

	//获取当前链接模块的链接ID
	GetConnID() uint32

	//获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr

	//发送数据，将数据发送给远程客户端
	//Send(data []byte) error、
	SendMsg(msgId uint32, data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error
