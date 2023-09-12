package AnalyzeServer

import (
	"MasterServer/Logs"
	"MasterServer/Minio"
	"encoding/json"
	"io/ioutil"
)

func InitServer() string {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	allclients = make(map[string]*ProfilerClient, 20) //暂定赋予20个解析客户端
	for i := 0; i < len(config.Client); i++ {
		allclients[config.Client[i].Ip] = &config.Client[i]
	}
	Minio.InitMinio(config.MinioServerPath, config.MinioBucket, config.MinioRawBucket)
	Logs.Loggers().Print("初始化服务器配置成功----")
	//测试是否反序列化成功
	// fmt.Print(config)
	serUrl := config.MasterServer.Ip + ":" + config.MasterServer.Port
	return serUrl
}

//初始化解析器，先ping一下
func InitAnalyzeClient() {
	for i := 0; i < len(config.Client); i++ {

	}
}
