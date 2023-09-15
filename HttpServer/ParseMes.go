package HttpServer

import (
	"MasterServer/Logs"
	"MasterServer/RabbitMqServer"
	"strings"
)

func StorageSucessParseMes(data string) {
	Logs.Loggers().Print("收到解析成功信号----" + data)
	dataPath := "./ServerQue/ParseQueSuccessQue"
	RabbitMqServer.PutData(dataPath, data)
}

func StorageAnalyzeParse(data string) {
	var uuid string
	if strings.Contains(data, "uuid") {
		uuid = GetUUID(data)
	}
	Logs.Loggers().Print("收到请求解析信号----" + data)
	dataPath := "./ServerQue/" + uuid + "_AnalyzeQue"
	RabbitMqServer.PutData(dataPath, data)
}

func GetUUID(mes string) string {
	var uuid string
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		if strings.Contains(str1[i], "uuid") { //解析uuid
			uid := strings.Split(str1[i], "=")
			uuid = uid[1]
			break
		}
	}
	return uuid
}
