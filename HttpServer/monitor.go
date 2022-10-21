package HttpServer

import (
	"UAutoServer/Logs"
	"fmt"
	"io"
	"net"
)

func ListenAndServer(address string) {
	//绑定监听地址
	listener, err := net.Listen("tcp", address)
	if err != nil {
		Logs.Print(fmt.Sprintf("listen err: %v", err))
	}
	Logs.Print("Http服务器启动成功！！！\n")
	defer listener.Close()

	for {
		//Accept 会一直阻塞直到有新的连接进来或者listen中断才会返回
		conn, err := listener.Accept()
		if err != nil {
			//通常由于listener被关闭无法继续监听导致的错误
			Logs.Print(fmt.Sprintf("accept err: %v", err))
		}
		//开启新的 goroutine处理该连接
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	Logs.Print("Http请求客户端连接成功----")
	//decoder := mahonia.NewDecoder("gbk")
	//decoder.NewReader()
	var resultData string
	for {
		tmp := make([]byte, 1024*1024)
		msg, err := conn.Read(tmp)
		if err != nil {
			//通常遇到的错误是连接中断或被关闭，用io.EOF表示
			if err == io.EOF {
				Logs.Print("connection close")
			} else {
				Logs.Print(err)
			}
			return
		}
		resultData += string(tmp[:msg])
		if tmp[msg-1] == '\n' {
			//DealReceivedMessage(resultData)
			resultData = ""
			b := []byte("Aceept Success.")
			conn.Write(b)
		}
		//fmt.Print(resultData, "已读取")
		//DealReceivedMessage(resultData)
		// b := []byte(msg)
		// //将收到的信息发送给客户端
		// conn.Write(b)
	}
}

// func DealReceivedMessage(msg string) {
// 	if strings.Contains(msg, "Start Sending Message") {
// 		beginMsg := strings.Split(msg, "|")[1]
// 		go ParseBegin(beginMsg)
// 	} else if strings.Contains(msg, "Stop Sending Message") {
// 		stopMsg := strings.Split(msg, "|")[1]
// 		go ParseStop(stopMsg)
// 	} else if strings.Contains(msg, "Give Up Sending Message") {
// 		giveupMsg := strings.Split(msg, "|")[1]
// 		go ParseGiveUp(giveupMsg)
// 		//todo delete datas
// 	} else if strings.Contains(msg, "End Sending Message") {
// 		endsendMsg := strings.Split(msg, "|")[1]
// 		go ParseEnd(endsendMsg)
// 	} else {
// 		go ParseReal(msg)
// 	}
// }
