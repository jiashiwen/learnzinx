package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块实现
*/

type Msghandle struct {
	//存放每个MsgID所对应的方法
	Apis map[uint32]ziface.IRouter

	//负责读取任务的消息队列
	TaskQueue []chan ziface.IRequest

	//业务工作的worker数量
	WorkerPoolSize uint32
}

//初始化/创建MsgHandler
func NewMsgHandle() *Msghandle {
	return &Msghandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *Msghandle) DoMsgHandler(request ziface.IRequest) {
	//从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is NOT bound ，need register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *Msghandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的Api处理方法师傅已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}

	//添加msg与Api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "succ!")

}

//启动一个worker工作池(开启工作池的动作只能发生一次
func (mh *Msghandle) StartWorkerPool() {
	//根据WorkerPoolSize分别开启Worker，每个Worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//启动一个worker

		//当前worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		//启动当前worker阻塞等待消息从channle传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])

	}

}

//启动一个worker工作流程
func (mh *Msghandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerId, "is started ...")

	//不断从阻塞等待对应消息队列消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)

		}
	}
}

//将消息交给TaskQueue处理
func (mh *Msghandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request MsgID=", request.GetMsgID(),
		" to WorkerID=", workerID)

	//将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request

}
