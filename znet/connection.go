package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前链接的Socket TCP套接字
	Conn *net.TCPConn

	//当前链接ID
	ConnID uint32

	//当前链接的状态
	isClosed bool

	//告知当前链接已经退出或停止 channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is Running")
	defer fmt.Println("ConnID=", c.ConnID, "Reader is exit,remote addr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		//读取客户端的数据到buf
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//得到当前conn数据的reuqest数据
		req := Request{
			conn: c,
			data: buf,
		}

		//从路由中，找到注册绑定Conn对应的router调用
		go func(request *Request) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

//启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start ... ConnID=", c.ConnID)

	//启动从当前链接读数据
	go c.StartReader()

	// TODO 启动从当前链接写数据
}

//停止链接
func (c *Connection) Stop() {
	fmt.Println("Conn Stopping ... ConnID=", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
