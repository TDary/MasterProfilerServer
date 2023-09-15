package AnalyzeServer

import "net"

//添加进入连接池
func AddConnectior(conn net.Conn) {
	var cPool ConnectPool
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	currentIp := remoteAddr.IP.String()
	isHas := false
	for _, val := range allconnector {
		if val.Ip == currentIp {
			isHas = true
			break
		}
	}
	cPool.Ip = currentIp
	cPool.Conn = conn
	if !isHas {
		allconnector = append(allconnector, cPool)
	}
}

//获取Conn连接池对象
func GetConn(ip string) net.Conn {
	for _, val := range allconnector {
		if val.Ip == ip {
			return val.Conn
		}
	}
	return nil
}

//清除一项连接
func ClearOneConn(ip string) {
	for key, val := range allconnector {
		if val.Ip == ip {
			val.Conn.Close()
			allconnector = append(allconnector[:key], allconnector[key+1:]...)
			break
		}
	}
}
