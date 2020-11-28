package znet

import "zinx/ziface"

//实现router时，先嵌入BaseRouter基类，根据需要对基类的方法进行重写
type BaseRouter struct {
}

//处理业务之前的
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

//处理conn之后的方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
