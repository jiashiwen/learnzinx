package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	//服务器名称
	Name string
	//ip版本
	IpVersion string
	//服务器ip
	Ip string
	//服务器端口
	Port int
	//当前server的消息管理模块，用来绑定MsgID对应的api关系
	MsgHandle ziface.IMsgHandle
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s,Listenner at IP:%s,Port:%d",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)

	fmt.Printf("[Zinx] Version:%s,MaxConn:%d,MaxPackageSize:%d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		//开启消息队列及worker工作池
		s.MsgHandle.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IpVersion, "err", err)
			return
		}

		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			dealConn := NewConnection(conn, cid, s.MsgHandle)
			cid++

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {

	//TODO 将服务器资源停止与回收

}

func (s *Server) Server() {
	//启动server
	s.Start()

	//阻塞
	select {}

}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("Add Router succ")
}

/*
初始化Server
*/

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IpVersion: "tcp4",
		Ip:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
	}

	return s
}
