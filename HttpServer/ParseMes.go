package HttpServer

import (
	"UAutoServer/Logs"
	"UAutoServer/RabbitMqServer"
)

func StopConnection() {

}

func StoragePaseMes(data string) {
	Logs.Print("收到请求解析信号----")
	RabbitMqServer.PutData("ParseQue", data)
}

func StorageSucessParseMes(data string) {
	Logs.Print("收到解析成功信号----")
	RabbitMqServer.PutData("SuccessQue", data)
}
