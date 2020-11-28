package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
模拟客户端
*/

func main() {
	fmt.Println("client start ...")
	time.Sleep(1 * time.Second)
	//链接服务器，得到一个conn
	conn, err := net.Dial("tcp", "127.0.0.1:7777")

	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	//调用Write 写数据
	for {
		//发送封包的msg消息
		dp := znet.NewDataPack()
		binaryMsg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx 0 Message")))

		if err != nil {
			fmt.Println("Pcak error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器回复一个message数据
		//读取流中Head部分得到ID和datalen

		//根据datalen二次读取，读data
		binaryhead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryhead); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		msgHead, err := dp.UnPack(binaryhead)
		if err != nil {
			fmt.Println("clent unpack mesghead error:", err)
			break
		}
		if msgHead.GetMegLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMegLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}
			fmt.Println("--> Recv server Msg: ID=", msg.Id, ",Len=", msg.DataLen, ",Data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)

	}
}
