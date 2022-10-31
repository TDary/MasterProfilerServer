package HttpServer

import (
	"UAutoServer/Logs"
	"UAutoServer/RabbitMqServer"
)

func StopConnection() {

}

func StorageParseMes(data string) {
	Logs.Print("收到请求解析信号----" + data)
	RabbitMqServer.PutData("/HttpServer/ParseQue", data)
}

func StorageSucessParseMes(data string) {
	Logs.Print("收到解析成功信号----" + data)
	RabbitMqServer.PutData("/HttpServer/ParseQueSuccessQue", data)
}

func GetStorageParseMes(data string) string {
	res := RabbitMqServer.GetData(data)
	Logs.Print("取出请求解析数据----" + res)
	return res
}
