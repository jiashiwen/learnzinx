package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//拆包封包单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/

	//创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listener err", err)
	}

	//创建go承载客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//拆包过程
				// 第一次从conn读取head
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
						return
					}
					if msgHead.GetMegLen() > 0 {
						//msg有数据，需要进行二次读取
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMegLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error:", err)
							return
						}
						//完成消息读取
						fmt.Println("--> Recv MsgID:", msg.Id, "DataLen:", msg.DataLen, "Data:", string(msg.Data))
					}
				}
				//第二次根据head读数据

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail err:", err)
		return
	}

	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送

	//封装第一个msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	//封装msg2
	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	//将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发送到客户端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
