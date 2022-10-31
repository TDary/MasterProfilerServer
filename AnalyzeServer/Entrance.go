package AnalyzeServer

import (
	"UAutoServer/Logs"
	"encoding/json"
	"io/ioutil"
)

type AllProfilerClient struct {
	Ip            string
	WorkerNumbers int
	WorkType      string
	State         bool
}

type MainTable struct {
	AppKey   string
	UUID     string
	RawFiles []string
}

type ConfigData struct {
	Client []AllProfilerClient
}

func Run() {
	InitClient()
	//HttpServer.ListenAndServer("10.11.144.31:8201")
}

func InitClient() {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var config ConfigData
	// var cData AllProfilerClient
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Error(err)
	}
	//测试是否反序列化成功
	// fmt.Print(config.Client[0].Ip)
}

//解析url进行结构化
func AnalyzeRequestUrl() {
	//从消息队列中取出解析的url进行操作
	//此处作为消费者,同时调用DataBase创建数据库表
	for true {

	}
}

func AddAnalyzeClient() {
	//当客户端解析器启动时会ping一次服务器，测试是否已将客户端解析器加入了组网,保证机器间正常运行
	//只有当ping通的情况下才会将开关打开
	//同时可以加入新的解析客户端
	//同时在此处进行启动客户端解析

}
