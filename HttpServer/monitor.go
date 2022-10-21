package HttpServer

import (
	"UAutoServer/Logs"
	"bytes"
	"fmt"
	"net"
	"net/url"
)

var AllProfilerClient map[string]string

func ListenAndServer(address string) {
	// tcp 连接，监听 8021 端口
	monitor, err := net.Listen("tcp", address)
	if err != nil {
		Logs.Error(err)
	}

	// 死循环，每当遇到连接时，调用 handle
	for {
		client, err := monitor.Accept()
		if err != nil {
			Logs.Print(err)
			continue
		}
		//每当有一个连接就使用协程
		go Handle(client)
	}
}

func Handle(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	Logs.Print(client.RemoteAddr())

	// 用来存放客户端数据的缓冲区
	var b [1024]byte
	//从客户端获取数据
	_, err := client.Read(b[:])
	if err != nil {
		Logs.Print(err)
		return
	}

	var method, URL string
	// 从客户端数据读入 method，url
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &URL)
	//获取到绝对url。ip端口号后面的内容
	hostPortURL, err := url.Parse(URL)
	if err != nil {
		Logs.Print(err)
		return
	}
	DealReceivedMessage(hostPortURL)
}

func DealReceivedMessage(msg *url.URL) {
	// if strings.Contains(msg, "Request Profiler Message") {
	// 	beginMsg := strings.Split(msg, "|")[1]
	// 	go StoragePaseMes(beginMsg)
	// } else if strings.Contains(msg, "Profiler Success Message") {
	// 	suce := strings.Split(msg, "|")[1]
	// 	go StorageSucessParseMes(suce)
	// } else if strings.Contains(msg, "Give Up Connect Message") {

	// } else {
	// }
}
