package main

import (
	"fmt"
	"net"
	"time"
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
		_, err := conn.Write([]byte("hello zinx v0.1"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err")
			return
		}

		fmt.Printf("server call back:\n %s,cnt=%d\n", buf, cnt)

		time.Sleep(1 * time.Second)

	}
}
