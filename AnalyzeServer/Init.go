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
	var client ClientState
	for _, val := range config.Client {
		client.Ip = val.Ip
		client.IpAddress = val.Ip + ":" + val.Port
		client.Num = val.WorkerNumbers
		client.State = "out"
		allAnalyzeClient = append(allAnalyzeClient, client)
	}
	InitAnalyzeClient()
	Minio.InitMinio(config.MinioServerPath, config.MinioBucket, config.MinioRawBucket)
	Logs.Loggers().Print("初始化服务器配置成功----")
	serUrl := config.MasterServer.Ip + ":" + config.MasterServer.Port
	return serUrl
}

//初始化解析器，先ping一下
func InitAnalyzeClient() {
	for _, val := range config.Client {
		address := val.Ip + ":" + val.Port
		rev := RequestClientState(address)
		if rev.State != "" {
			for key, val2 := range allAnalyzeClient {
				if val2.IpAddress == address {
					allAnalyzeClient[key].State = rev.State
					break
				}
			}
		}
	}
}
