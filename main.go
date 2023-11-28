//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"MasterServer/AnalyzeServer"
	"MasterServer/HttpServer"
	"MasterServer/Logs"
)

func main() {
	//日志初始化
	Logs.Init()
	Logs.Loggers().Print("Welcome to use ServerMaster")
	//初始化服务器配置
	ServerUrl := AnalyzeServer.InitServer()
	//启动开始处理完成解析消息系统
	go AnalyzeServer.AnalyzeSuccessToMerge()
	//检测失败任务
	go AnalyzeServer.CheckFailedAnalyzeData()
	//启动socket监听
	HttpServer.ListenAndServer(ServerUrl)
}

//todo:宕机重启的情况下去检查一下数据库，然后根据各个状态来分配，有未解析的完的就去检查一下存储系统是否有解析完的数据
