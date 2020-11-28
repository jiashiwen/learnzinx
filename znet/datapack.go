package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

//拆包分包实例初始化
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + Id U你32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建存放byte的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将DataLen写入Buff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMegLen()); err != nil {
		return nil, err
	}

	//将MsgId写入Buff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写入Buff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法,读出head信息，再根据head中data长度再读内容
func (dp *DataPack) UnPack(binarryData []byte) (ziface.IMessage, error) {
	//创建从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binarryData)
	//解析Head信息，得到datalen和MsgId
	msg := &Message{}

	//读Datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("Too large msg data recv!")
	}

	return msg, nil
}
