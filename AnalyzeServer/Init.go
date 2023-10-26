package AnalyzeServer

import (
	"MasterServer/Logs"
	"MasterServer/Minio"
	"encoding/json"
	"io/ioutil"
	"os"
)

func InitServer() string {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	_, err = os.Stat(config.MergePath)
	if err != nil {
		Logs.Loggers().Printf("当前文件夹%s不存在，重新创建中！", config.MergePath)
		os.Mkdir(config.MergePath, 0755)
	}

	filepath := "./ServerQue"
	_, err = os.Stat(filepath)
	if err != nil {
		os.Mkdir(filepath, 0755)
	}

	for _, val := range config.Client {
		var client ClientState
		client.Ip = val.Ip
		client.IpAddress = val.Ip + ":" + val.Port
		client.Num = val.WorkerNumbers
		client.State = "out"
		allAnalyzeClient = append(allAnalyzeClient, client)
	}

	Minio.InitMinio(config.MinioServerPath, config.MinioBucket, config.MinioRawBucket)
	serUrl := config.MasterServer.Ip + ":" + config.MasterServer.Port
	Logs.Loggers().Print("初始化服务器配置成功----")
	return serUrl
}
