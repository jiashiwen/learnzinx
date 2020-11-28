package ziface

/*
封包、拆包模块
处理TCP连接的数据流 ，用于处理TCP粘包问题
*/

type IDataPack interface {
	//获取包头长度
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte, error)
	//拆包方法
	UnPack([]byte) (IMessage, error)
}
