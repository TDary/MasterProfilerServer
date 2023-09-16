package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/RabbitMqServer"
	"strconv"
	"strings"
)

//停止采集信号
func StopAnalyzeRequest(mes string) {
	var data EndData
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		if strings.Contains(str1[i], "uuid") { //解析gameid
			uid := strings.Split(str1[i], "=")
			data.UUID = uid[1]
		} else if strings.Contains(str1[i], "lastfile") { //解析rawfiles
			file := strings.Split(str1[i], "=")
			data.LastRawFile = file[1]
		} else {
			Logs.Loggers().Print("不存在对额外参数的解析:" + str1[i])
		}
	}
	stopMsg = append(stopMsg, data)
}

//获取空闲状态的解析器,供发送请求用
func GetIdleAnalyzeClient() {
	for _, val := range config.Client {
		address := val.Ip + ":" + val.Port
		res := RequestClientState(address)
		if res.State == "idle" {
			var current ClientState
			current.IpAddress = address
			current.State = res.State
			current.Num = res.Num
			ishas := false
			for i := 0; i < len(allAnalyzeClient); i++ {
				if allAnalyzeClient[i].IpAddress == address {
					allAnalyzeClient[i].State = res.State
					allAnalyzeClient[i].Num = res.Num
					ishas = true
					break
				}
			}
			if !ishas {
				allAnalyzeClient = append(allAnalyzeClient, current)
				break
			}
		}
	}

}

//检测最后一份源文件
func GetLastRawFileIsSend(uuid string) int {
	for _, val := range stopMsg {
		if val.UUID == uuid {
			return 1

		}
	}
	return -1 //继续等待
}

//将原始文件进行排序
func SortRawFils(rawfiles []string) {
	for i := 0; i < len(rawfiles)-1; i++ {
		for j := 0; j < len(rawfiles)-i-1; j++ {
			fnum1, err := strconv.Atoi(strings.Split(rawfiles[j], ".")[0])
			if err != nil {
				Logs.Loggers().Print("转换失败----", err.Error())
				return
			}
			fnum2, err := strconv.Atoi(strings.Split(rawfiles[j+1], ".")[0])
			if err != nil {
				Logs.Loggers().Print("转换失败----", err.Error())
				return
			}
			if fnum1 > fnum2 {
				rawfiles[j], rawfiles[j+1] = rawfiles[j+1], rawfiles[j]
			}
		}
	}
}

//发送真正的解析请求
func AnalyzeBegin(analze string, databaseData string) {
	AddOneForSubTable(databaseData) //添加数据库子任务表
	//res := GetAnalyzeData(analze)
	//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
	// GetIdleAnalyzeClient()
	for _, val := range allAnalyzeClient {
		n, err := GetConn(val.Ip, "anaclient").Write([]byte(analze))
		if err != nil && n == 0 {
			Logs.Loggers().Print("发送解析消息失败----", err.Error())
			break
		} else {
			//
			Logs.Loggers().Print("发送长度：", n)
			break
		}
	}
}

//解析url进行结构化并创建数据库表数据
func AnalyzeRequest(data string) {
	//此处作为消费者,同时调用DataBase创建数据库表
	mtable := ReceiveMes(data)
	var rawFiles []string
	quePath := "./ServerQue/" + mtable.UUID + "_AnalyzeQue"
	for {
		state := GetLastRawFileIsSend(mtable.UUID)
		if state == 1 {
			//队列中拿rawfilename
			for {
				getanalyzeData := RabbitMqServer.GetData(quePath)
				if getanalyzeData != "" {
					res := GetAnalyzeData(getanalyzeData)
					rawFiles = append(rawFiles, res.RawFile)
				} else {
					//更新数据库主表，先排序，这里用了冒泡
					SortRawFils(rawFiles)
					DataBase.UpdateMainTable(mtable.AppKey, mtable.UUID, rawFiles) //更新源文件队列
					if isMergeStop {
						isMergeStop = false //通知开始进行状态修改并合并
					}
					break
				}
			}
			break
		} else {
			//waiting
		}
	}
}

//查询正在运行的工作机
func CheckKey(key string) bool {
	for _, val := range allAnalyzeClient {
		if val.Ip == key {
			return true
		}
	}
	return false
}

//重新解析
func ReProfilerAna(data string) {
	spldata := strings.Split(data, "&")
	var uuid string
	var rawfile string
	for i := 0; i < len(spldata); i++ {
		if strings.Contains(spldata[i], "uuid") {
			uid := strings.Split(spldata[i], "=")
			uuid = uid[1]
		} else if strings.Contains(spldata[i], "rawfile") {
			file := strings.Split(spldata[i], "=")
			rawfile = file[1]
		}
	}
	DataBase.FindAndModify(uuid, rawfile)
}

//添加一项子表任务
func AddOneForSubTable(data string) {
	var subt DataBase.SubTable
	spldata := strings.Split(data, "&")
	for i := 0; i < len(spldata); i++ {
		if strings.Contains(spldata[i], "uuid") {
			uid := strings.Split(spldata[i], "=")
			subt.UUID = uid[1]
		} else if strings.Contains(spldata[i], "rawfile") {
			file := strings.Split(spldata[i], "=")
			subt.RawFile = file[1]
		} else if strings.Contains(spldata[i], "ip") {
			ip := strings.Split(spldata[i], "=")
			subt.AnalyzeIP = ip[1]
		}
	}
	subt.State = 0
	InsertSubTableBySub(subt) //插入一条子任务
}