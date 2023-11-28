package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
	"MasterServer/Tools"
	"encoding/json"
	"io/ioutil"
	"os"
)

func InitServer() string {
	var data, _ = ioutil.ReadFile("./ServerConfig.dat")
	key := []byte("eb3386a8a8f57a579c93fdfb33ec9471") // 加密密钥，长度为16, 24, 或 32字节，对应AES-128, AES-192, AES-256
	decryptedData, err := Tools.Decrypt(data, key)
	if err != nil {
		Logs.Loggers().Print(err)
		return ""
	}
	err = json.Unmarshal(decryptedData, &config)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	_, err = os.Stat(config.Minioconfig.MergePath)
	if err != nil {
		Logs.Loggers().Printf("当前文件夹%s不存在，重新创建中！", config.Minioconfig.MergePath)
		os.Mkdir(config.Minioconfig.MergePath, 0755)
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
		client.AnalyzeType = val.WorkType
		allAnalyzeClient = append(allAnalyzeClient, client)
	}
	//初始化数据库配置与Minio服务配置
	DataBase.InitDB(config.Database.Address, config.Database.DBName, config.Database.Collection.MainTable, config.Database.Collection.SubTable, config.Database.Collection.FunRow, config.Database.Collection.SimpleTable, config.Database.Collection.FunPath)
	Minio.InitMinio(config.Minioconfig.MinioServerPath, config.Minioconfig.MinioBucket, config.Minioconfig.MinioRawBucket, config.Minioconfig.UserName, config.Minioconfig.PassWord)
	serUrl := config.MasterServer.Ip + ":" + config.MasterServer.Port
	Logs.Loggers().Print("初始化服务器配置成功----")
	return serUrl
}
