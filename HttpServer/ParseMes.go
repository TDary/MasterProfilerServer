package HttpServer

import (
	"UAutoServer/AnalyzeServer"
	"UAutoServer/Logs"
	"UAutoServer/RabbitMqServer"
)

func StorageParseMes(data string) {
	Logs.Loggers().Print("收到请求解析信号----" + data)
	RabbitMqServer.PutData("/HttpServer/ParseQue", data)
	AnalyzeServer.ChangeMessage()
}

func StorageSucessParseMes(data string) {
	Logs.Loggers().Print("收到解析成功信号----" + data)
	RabbitMqServer.PutData("/HttpServer/ParseQueSuccessQue", data)
	AnalyzeServer.ChangeSuccessMessage()
}
