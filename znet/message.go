package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

//创建Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

//获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息长度
func (m *Message) GetMegLen() uint32 {
	return m.DataLen
}

//获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息的ID
func (m *Message) SetMsgID(id uint32) {
	m.Id = id
}

//设置消息的长度
func (m *Message) SetMsgLen(length uint32) {
	m.DataLen = length
}

//设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
