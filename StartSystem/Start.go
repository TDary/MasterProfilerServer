package StartSystem

import (
	"UAutoServer/AnalyzeServer"
	"UAutoServer/DataBase"
	"UAutoServer/HttpServer"
)

func Run() {
	//初始化数据库配置
	DataBase.InitDB()
	//初始化服务器配置
	AnalyzeServer.InitClient()
	//启动解析url系统
	go AnalyzeServer.AnalyzeRequestUrl()
	//启动Http监听
	HttpServer.ListenAndServer("10.11.144.31:8201")
}
