package AnalyzeServer

import (
	"MasterServer/Logs"
	"encoding/json"
	"io/ioutil"
)

var config ConfigData
var isStop bool                           //请求解析处理控制信号
var isAnalyzeStop bool                    //完成解析处理控制信号
var allclients map[string]*ProfilerClient //解析客户端及服务端配置

type SuccessData struct {
	UUID    string
	IP      string
	RawFile string
}

type ProfilerClient struct {
	Ip            string
	Port          string
	WorkerNumbers int
	WorkType      string
	State         bool
}

type MergeServerConfig struct {
	Ip   string
	Port string
}
type ConfigData struct {
	Client      []ProfilerClient
	MergeServer MergeServerConfig
}

func InitClient() {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	allclients = make(map[string]*ProfilerClient, 20) //暂定赋予20个解析客户端
	for i := 0; i < len(config.Client); i++ {
		allclients[config.Client[i].Ip] = &config.Client[i]
	}
	Logs.Loggers().Print("初始化服务器配置成功----")
	//测试是否反序列化成功
	// fmt.Print(config)
}

func ChangeMessage() {
	if isStop {
		isStop = false
	} else {
		return
	}
}

func ChangeSuccessMessage() {
	if isStop {
		isAnalyzeStop = false
	} else {
		return
	}
}
