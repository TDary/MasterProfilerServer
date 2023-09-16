package AnalyzeServer

import "net"

//添加进入连接池
func AddConnectior(conn net.Conn, machine string) {
	var cPool ConnectPool
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	currentIp := remoteAddr.IP.String()
	isHas := false
	for _, val := range allconnector {
		if val.Ip == currentIp && val.Marchine == machine {
			isHas = true
			val.Conn = conn
			break
		}
	}
	cPool.Ip = currentIp
	cPool.Conn = conn
	cPool.Marchine = machine
	if !isHas {
		allconnector = append(allconnector, cPool)
	}
}

//获取Conn连接池对象
func GetConn(ip string, machine string) net.Conn {
	for _, val := range allconnector {
		if val.Ip == ip && machine == val.Marchine {
			return val.Conn
		}
	}
	return nil
}
