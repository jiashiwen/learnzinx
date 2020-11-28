package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//hellozinx 自定义路由
type HelloZinx struct {
	znet.BaseRouter
}

//test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")
	//读取客户端数据，再回写ping..ping..ping..
	fmt.Println("recv from client:msgID=", request.GetMsgID(), ",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping..ping..ping.."))
	if err != nil {
		fmt.Println(err)
	}
}

//hellozinx Handle
func (h *HelloZinx) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")
	//读取客户端数据，再回写ping..ping..ping..
	fmt.Println("recv from client:msgID=", request.GetMsgID(), ",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("Hello zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer("myzinx")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinx{})
	s.Server()
}
