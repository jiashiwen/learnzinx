package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前链接的Socket TCP套接字
	Conn *net.TCPConn

	//当前链接ID
	ConnID uint32

	//当前链接的状态
	isClosed bool

	//告知当前链接已经退出或停止 channel(由Reader告知Writer客户端退出)
	ExitChan chan bool

	//无缓冲的管道，用于读、写Goroutine之间的消息
	msgChan chan []byte

	//消息管理的MsgID和对应的业务处理api关系
	MsgHandler ziface.IMsgHandle
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is Running]")
	defer fmt.Println("ConnID=", c.ConnID, "[Reader is exit,remote addr is] ", c.RemoteAddr())
	defer c.Stop()

	for {

		//创建拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}

		//拆包，得到MsgID和MsgDataLen
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据datalen再次读取Data，放入msg.Data
		var data []byte
		if msg.GetMegLen() > 0 {
			data = make([]byte, msg.GetMegLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Print("read msg data error", err)
				break
			}
		}

		//得到当前conn数据的reuqest数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//根据绑定好的MsgId找到对应的api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

//写消息的Goroutine，专门用于发送消息给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Write Goroutine start ...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	//不断阻塞的等待channel的消息，进行回写客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

//SendMsg方法将数据先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closee when send msg")
	}

	//将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Print("Pack error msg id= ", msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

//启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start ... ConnID=", c.ConnID)

	//启动从当前链接读数据
	go c.StartReader()

	// 启动从当前链接写数据goroutine
	go c.StartWriter()
}

//停止链接
func (c *Connection) Stop() {
	fmt.Println("Conn Stopping ... ConnID=", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)

}

//获取当前链接的绑定socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
