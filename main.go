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

	// test := "{\"code\":200,\"state\":\"idle\",\"num\":4}"
	// var rece AnalyzeServer.ReceiveDate
	// err := json.Unmarshal([]byte(test), &rece)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Print(rece)
}

//todo:
//服务器被强行关机的情况要做处理!!!!!!!断线重连的情况等
//宕机重启的情况下去检查一下数据库，然后根据各个状态来分配，有未解析的完的就去检查一下存储系统是否有解析完的数据
//新增清除废弃数据功能，每到凌晨触发，删除已取消采集的任务

//反序列化测试  成功
// func Test() {
// 	rawDataPath := "D:\\000\\result.bin"
// 	bytedata, err := ioutil.ReadFile(rawDataPath)
// 	if err != nil {
// 		//打开失败
// 		fmt.Print("打开分析文件失败----", rawDataPath)
// 		return
// 	}
// 	currentSimpleData := &Data.AllCaseFunRow{}
// 	err = proto.Unmarshal(bytedata, currentSimpleData)
// 	if err != nil {
// 		Logs.Loggers().Print("反序列化失败----", err.Error())
// 		return
// 	}
// 	fmt.Print(currentSimpleData.Allvalues[1].Frames[0].Gcalloc)
// }
