package ziface

/*
将请求的消息封装到Message中
*/

type IMessage interface {
	//获取消息的ID
	GetMsgId() uint32

	//获取消息长度
	GetMegLen() uint32

	//获取消息内容
	GetData() []byte

	//设置消息的ID
	SetMsgID(id uint32)
	//设置消息的长度
	SetMsgLen(length uint32)
	//设置消息的内容
	SetData(data []byte)
}
