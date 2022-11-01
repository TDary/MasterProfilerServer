package AnalyzeServer

import (
	"UAutoServer/DataBase"
	"UAutoServer/Logs"
	"UAutoServer/Tools"
	"strconv"
	"strings"
)

//解析url进行结构化
func AnalyzeRequestUrl() {
	//从消息队列中取出解析的url进行操作
	//此处作为消费者,同时调用DataBase创建数据库表
	var getUrlData string
	isStop = false
	for true {
		if !isStop {
			getUrlData = GetStorageParseMes("/HttpServer/ParseQue")
			if getUrlData != "" {
				DataBase.ReceiveMes(getUrlData)
			} else {
				Logs.Loggers().Print("队列已空，进入阻塞状态...") //使用通道式消息，知道接收到有解析消息才会解开
				isStop = true
			}
		}
	}
}

//对解析成功的消息进行检查判断是否可以进行合并操作
func AnalyzeSuccessUrl() {

}

func AddAnalyzeClient(data string) {
	//当客户端解析器启动时会ping一次服务器，测试是否已将客户端解析器加入了组网,保证机器间正常运行
	//只有当ping通的情况下才会将开关打开
	//同时可以加入新的解析客户端
	//同时在此处进行启动客户端解析
	nowstr := strings.Split(data, "&")
	nowip := strings.Split(nowstr[0], "=")
	if Tools.GetKey(nowip[1], allclientIP) {
		for i := 0; i < len(config.Client); i++ {
			if config.Client[i].Ip == nowip[1] {
				config.Client[i].State = true
				Logs.Loggers().Print("该客户端解析打开成功----IP:" + nowip[1])
				return
			}
		}
	} else {
		var newClient ProfilerClient
		num := strings.Split(nowstr[1], "=")
		newClient.Ip = nowip[1]
		nums, err := strconv.Atoi(num[1])
		if err != nil {
			Logs.Loggers().Print("转换类型失败----")
			return
		}
		newClient.WorkerNumbers = nums
		newClient.WorkType = "Analyze"
		newClient.State = true
		config.Client = append(config.Client, newClient)
		Logs.Loggers().Print("识别到新解析客户端，加入组网成功----")
	}
}
