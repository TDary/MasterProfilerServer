package AnalyzeServer

import (
	"UAutoServer/Logs"
	"encoding/json"
	"io/ioutil"
)

var config ConfigData
var allclientIP map[int]string
var isStop bool

type ProfilerClient struct {
	Ip            string
	WorkerNumbers int
	WorkType      string
	State         bool
}

type ConfigData struct {
	Client []ProfilerClient
}

func InitClient() {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	allclientIP = make(map[int]string, 20) //暂定赋予20个解析客户端
	for i := 0; i < len(config.Client); i++ {
		allclientIP[i] = config.Client[i].Ip
	}
	Logs.Loggers().Print("初始化服务器配置成功----")
	//测试是否反序列化成功
	// fmt.Print(config.Client[0].Ip)
}

func ChangeMessage() {
	if isStop {
		isStop = false
	} else {
		return
	}
}
