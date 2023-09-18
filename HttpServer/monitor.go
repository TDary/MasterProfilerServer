package HttpServer

import (
	"MasterServer/AnalyzeServer"
	"MasterServer/Logs"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// 处理连接逻辑
	// 在这里可以读取和写入数据

	message := "Welcome to the server!\n"
	conn.Write([]byte(message))

	for {
		// 示例：从客户端读取数据并打印
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err != nil && n == 0 {
			Logs.Loggers().Printf("Error reading from connection: %s", err.Error())
			//断开连接，清除池子
			return
		}
		if len(buffer) != 0 {
			res := string(buffer[:n])
			if strings.Contains(res, "startanalyze") {
				Logs.Loggers().Print("接收到开始采集消息----", res)
				beginMsg := strings.Split(res, "?")[1] //startanalyze?...
				go AnalyzeServer.AnalyzeRequest(beginMsg)
				message = "ok"
				conn.Write([]byte(message))
			} else if strings.Contains(res, "requestanalyze") {
				Logs.Loggers().Print("接收到申请解析源文件的消息----", res)
				ana := strings.Split(res, "?")[1] //requestanalyze?uuid=test&rawfile=123123.zip&rawfilename=uuid/1231.zip&unityversion=12313&analyzebucket=ads&analyzeType=
				go StorageAnalyzeParse(ana)
				go AnalyzeServer.AnalyzeBegin(res, ana)
				message = "ok"
				conn.Write([]byte(message))
			} else if strings.Contains(res, "successprofiler") {
				Logs.Loggers().Print("接收到解析成功消息----", res)
				suce := strings.Split(res, "?")[1]
				go AnalyzeServer.ParseSuccessData(suce)
				message = "ok"
				conn.Write([]byte(message))
			} else if strings.Contains(res, "rquestclient") {
				Logs.Loggers().Print("接收到解析器请求消息----", res)
				message = "ok"
				conn.Write([]byte(message))
			} else if strings.Contains(res, "stopanalyze") {
				Logs.Loggers().Print("接收到停止采集消息----", res)
				stopMsg := strings.Split(res, "?")[1]
				go AnalyzeServer.StopAnalyzeRequest(stopMsg)
				message = "ok"
				conn.Write([]byte(message))
			} else if strings.Contains(res, "ReAnalyze") {
				Logs.Loggers().Print("接收到重新解析消息----", res)
				req := strings.Split(res, "?")[1]
				go AnalyzeServer.ReProfilerAna(req)
			} else if strings.Contains(res, "markeid") {
				Logs.Loggers().Print("接收到加入连接消息----", res)
				req := strings.Split(res, "?")[1]
				AnalyzeServer.AddConnectior(conn, req)
			} else {
				Logs.Loggers().Print("receive Data:", res)
			}
		}
	}
}

func ListenAndServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		Logs.Loggers().Fatalf("Failed to listen: %s", err.Error())
	}
	defer listener.Close()
	Logs.Loggers().Printf("Server listening on %s\n", address)
	for {
		// 接受新连接
		conn, err := listener.Accept()
		if err != nil {
			Logs.Loggers().Printf("Failed to accept connection: %s", err.Error())
			continue
		}

		// 处理连接
		go HandleConnection(conn)
	}
}
