package AnalyzeServer

import (
	"MasterServer/Logs"
	"net"
)

// 添加进入连接池
func AddConnectior(conn net.Conn, machine string) {
	var cPool ConnectPool
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	currentIp := remoteAddr.IP.String()
	isHas := false
	for k, val := range allconnector {
		if val.Ip == currentIp && val.Marchine == machine {
			isHas = true
			allconnector[k].Conn = conn
			break
		}
	}
	cPool.Ip = currentIp
	cPool.Conn = conn
	cPool.Marchine = machine
	if !isHas {
		allconnector = append(allconnector, cPool)
	}
	Logs.Loggers().Print("连接池数量", len(allconnector))
	//将解析器状态调整为空闲，设置之后可以进行解析操作
	for key, val := range allAnalyzeClient {
		if val.Ip == currentIp && machine == "anaclient" {
			allAnalyzeClient[key].State = "idle"
		}
	}
}

// 获取Conn连接池对象
func GetConn(ip string, machine string) net.Conn {
	for _, val := range allconnector {
		if val.Ip == ip && machine == val.Marchine {
			return val.Conn
		}
	}
	return nil
}

// 断开连接，去除连接池中的连接
func CloseConnect(ip string, machine string) {
	for _, val := range allconnector {
		if val.Ip == ip && machine == val.Marchine {
			val.Conn = nil
		}
	}
}
