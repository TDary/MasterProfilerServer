package AnalyzeServer

import (
	"UAutoServer/Logs"
	"UAutoServer/RabbitMqServer"
)

func GetStorageParseMes(data string) string {
	res := RabbitMqServer.GetData(data)
	Logs.Loggers().Print("取出请求解析数据----" + res)
	return res
}
