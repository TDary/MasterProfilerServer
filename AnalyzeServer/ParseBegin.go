package AnalyzeServer

import (
	"UAutoServer/DataBase"
	"UAutoServer/Logs"
	"time"
)

func ParseEntrance() {
	for true {
		CheckSubTable()
		time.Sleep(5 * time.Minute)
	}
}

func SendBeginMessage() {
	//发送开始解析的相关数据信息
	//前提是已经创建好数据表
}

func CheckSubTable() {
	currentState := 0
	andata := DataBase.FindSTbyState(currentState)
	for _, val := range andata {
		Logs.Loggers().Print(val)
	}
}

func ReduceRunC(ip string, rcount int) {
	if CheckKey(ip) {
		if allclients[ip].WorkerNumbers == 0 {
			allclients[ip].State = false
			Logs.Loggers().Print("IP:" + ip + "的可用机器已用完，正在等待释放，该机器强行设置进入关闭状态。")
		} else {
			allclients[ip].WorkerNumbers = allclients[ip].WorkerNumbers - rcount
		}
	}
}

func GetTotalRunC() int {
	var totalCount int
	for _, value := range allclients {
		totalCount += value.WorkerNumbers
	}
	return totalCount
}
