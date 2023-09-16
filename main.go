//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"MasterServer/AnalyzeServer"
	"MasterServer/DataBase"
	"MasterServer/HttpServer"
	"MasterServer/Logs"
)

func main() {
	Logs.Loggers().Print("Welcome to use ServerMaster")
	//初始化数据库配置
	DataBase.InitDB()
	//初始化服务器配置;
	ServerUrl := AnalyzeServer.InitServer()
	//启动开始处理完成解析消息系统
	go AnalyzeServer.AnalyzeSuccessUrl()
	//启动socket监听
	HttpServer.ListenAndServer(ServerUrl)
}

//todo:
//服务器被强行关机的情况要做处理!!!!!!!断线重连的情况等
//宕机重启的情况下去检查一下数据库，然后根据各个状态来分配，有未解析的完的就去检查一下存储系统是否有解析完的数据
//新增清除废弃数据功能，每到凌晨触发，删除已取消采集的任务
