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

//test PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle ...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("befor handele \n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}

}

//test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("handele \n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

//test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router Post Handle ...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("after handele\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("myzinx")
	s.AddRouter(&PingRouter{})
	s.Server()
}
