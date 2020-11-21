package ziface

type IRouter interface {
	//处理业务之前的Hook
	PreHandle(request IRequest)

	//处理conn业务的主方法Hook
	Handle(request IRequest)

	//处理conn之后的方法
	PostHandle(request IRequest)
}
