package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
	"MasterServer/Tools"
	"encoding/json"
	"os"
)

func InitServer() string {
	failedquePath = "./ServerQue/" + "FailedAnalyzeQue"
	var data, _ = os.ReadFile("./ServerConfig.dat")
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
		os.MkdirAll(config.Minioconfig.MergePath, 0755)
	}
	switch config.AnalyzeMode {
	case "local":
		Logs.Loggers().Print("当前解析模式为本地单机模式，需要启动MasterClient客户端")
	case "distributed":
		Logs.Loggers().Print("当前解析模式为分布式联网模式。")
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

// // 本地解析模式下启用 且程序配置正确位置
// func StartMasterClient() {
// 	time.Sleep(5 * time.Second) //等一会 主服务器优先启动
// 	Logs.Loggers().Print("Start MasterClient server process----")
// 	cmd := exec.Command("./MasterClient.exe")
// 	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
// 	_, err := cmd.CombinedOutput()
// 	if err != nil {
// 		Logs.Loggers().Fatal("Failed to start MasterClient server.", err.Error())
// 		return
// 	}
// }
