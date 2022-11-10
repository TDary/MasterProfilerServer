package AnalyzeServer

import (
	"MasterServer/Logs"
	"encoding/json"
	"io/ioutil"
)

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
