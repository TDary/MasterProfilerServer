package StartSystem

import (
	"MasterServer/AnalyzeServer"
	"MasterServer/DataBase"
	"MasterServer/HttpServer"
)

func Run() {
	//初始化数据库配置
	DataBase.InitDB()
	//初始化服务器配置
	ServerUrl := AnalyzeServer.InitClient()
	//启动解析url系统
	go AnalyzeServer.AnalyzeRequestUrl()
	// //启动开始解析消息系统
	go AnalyzeServer.ParseEntrance()
	// //启动开始处理完成解析消息系统
	go AnalyzeServer.AnalyzeSuccessUrl()
	// //启动Http监听
	HttpServer.ListenAndServer(ServerUrl)
}
