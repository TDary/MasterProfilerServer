package AnalyzeServer

import (
	"UAutoServer/DataBase"
	"UAutoServer/Logs"
	"UAutoServer/Tools"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

var config ConfigData
var allclientIP map[int]string

type ProfilerClient struct {
	Ip            string
	WorkerNumbers int
	WorkType      string
	State         bool
}

type ConfigData struct {
	Client []ProfilerClient
}

func InitClient() {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Error(err)
	}
	allclientIP = make(map[int]string, 20) //暂定赋予20个解析客户端
	for i := 0; i < len(config.Client); i++ {
		allclientIP[i] = config.Client[i].Ip
	}
	Logs.Print("初始化服务器配置成功----")
	//测试是否反序列化成功
	// fmt.Print(config.Client[0].Ip)
}

//解析url进行结构化
func AnalyzeRequestUrl() {
	//从消息队列中取出解析的url进行操作
	//此处作为消费者,同时调用DataBase创建数据库表
	var getUrlData string
	var isStop bool
	isStop = false
	for true {
		if !isStop {
			getUrlData = GetStorageParseMes("/HttpServer/ParseQue")
			if getUrlData != "" {
				DataBase.ReceiveMes(getUrlData)
			} else {
				Logs.Print("队列已空，进入阻塞状态...") //使用通道式消息，知道接收到有解析消息才会解开
				isStop = true
			}
		}
	}
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
				Logs.Print("该客户端解析打开成功----IP:" + nowip[1])
				return
			}
		}
	} else {
		var newClient ProfilerClient
		num := strings.Split(nowstr[1], "=")
		newClient.Ip = nowip[1]
		nums, err := strconv.Atoi(num[1])
		if err != nil {
			Logs.Print("转换类型失败----")
			return
		}
		newClient.WorkerNumbers = nums
		newClient.WorkType = "Analyze"
		newClient.State = true
		config.Client = append(config.Client, newClient)
		Logs.Print("识别到新解析客户端，加入组网成功----")
	}
}
