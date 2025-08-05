package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/RabbitMqServer"
	"os"
	"strconv"
	"strings"
	"time"
)

// 停止采集信号
func StopGatherRequest(mes string) {
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

// 检测最后一份源文件
func GetLastRawFileIsSend(uuid string) int {
	for _, val := range stopMsg {
		if val.UUID == uuid {
			return 1
		}
	}
	return -1 //继续等待
}

// 将原始文件进行排序
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

// 发送开始采集失败的消息
func SendFailToGather(uuid string, ip string) {
	msg := "开始采集失败，当前存在重复的UUID" + uuid
	n, err := GetConn(ip, "collector").Write([]byte(msg))
	if err != nil && n == 0 {
		Logs.Loggers().Print("发送消息失败----", err.Error())
	}
}

// 发送真正的解析请求
func AnalyzeBegin(analze string, databaseData string) {
	analyzetype := AddOneForSubTable(databaseData) //添加数据库子任务表
	//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
	for _, val := range allAnalyzeClient {
		if val.State == "idle" && val.AnalyzeType == analyzetype {
			n, err := GetConn(val.Ip, "anaclient").Write([]byte(analze))
			if err != nil && n == 0 {
				Logs.Loggers().Print("发送解析消息失败----", err.Error())
				break
			} else {
				//Logs.Loggers().Print("发送请求解析成功~")
				// Logs.Loggers().Print("发送长度：", n)
				break
			}
		}
	}
}

// 处理开始采集以及采集过程消息，并在最后更新源文件列表
func StartGatherRequest(data string) {
	//此处作为消费者,同时调用DataBase创建数据库表
	mtable := ReceiveMes(data)
	if mtable.UUID == "" {
		Logs.Loggers().Print("由于插入数据库表数据失败，退出当前执行.")
		//发送采集失败消息进行停止
		SendFailToGather(mtable.UUID, mtable.CollectorIp)
		//断开连接
		CloseConnect(mtable.CollectorIp, "collector")
		return
	}
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
					DataBase.UpdateMainTable(mtable.AppKey, mtable.UUID, rawFiles) //更新源文件队列,+合并状态为可合并3
					break
				}
			}
			os.Remove(quePath) //拿完队列，将文件删除
			break
		} else {
			//waiting
		}
	}
}

// 查询正在运行的工作机
func CheckKey(key string) bool {
	for _, val := range allAnalyzeClient {
		if val.Ip == key {
			return true
		}
	}
	return false
}

// 重新解析
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
	currentTime := time.Now()
	unixTime := currentTime.Unix()
	DataBase.FindAndModify(uuid, rawfile, 0, unixTime) //修改任务状态
	//发送解析请求
	for _, val := range allAnalyzeClient {
		if val.State == "idle" {
			n, err := GetConn(val.Ip, "anaclient").Write([]byte(data))
			if err != nil && n == 0 {
				Logs.Loggers().Print("发送解析消息失败----", err.Error())
				break
			} else {
				//
				// Logs.Loggers().Print("发送长度：", n)
				break
			}
		}
	}
}

// 添加一项子表任务
func AddOneForSubTable(data string) string { //返回文件解析类型
	var subt DataBase.SubTable
	var anaType string
	spldata := strings.Split(data, "&")
	for i := 0; i < len(spldata); i++ {
		if strings.Contains(spldata[i], "uuid") {
			uid := strings.Split(spldata[i], "=")
			subt.UUID = uid[1]
		} else if strings.Contains(spldata[i], "rawfile") {
			file := strings.Split(spldata[i], "=")
			if len(file) < 2 {
				// 如果没有rawfile,说明不是用raw文件解析
			} else {
				subt.RawFile = file[1]
			}
		} else if strings.Contains(spldata[i], "snapfile") {
			file := strings.Split(spldata[i], "=")
			if len(file) < 2 {
				// 如果没有snapfile,说明不是用snap快照文件解析
			} else {
				subt.SnapFile = file[1]
			}
		} else if strings.Contains(spldata[i], "analyzetype") {
			anatype := strings.Split(spldata[i], "=")
			anaType = anatype[1]
		}
	}
	subt.AnalyzeIP = ""
	subt.State = 0
	currentTime := time.Now()
	beginUnixTime := currentTime.Unix()
	subt.AnalyzeBegin = beginUnixTime
	InsertSubTableBySub(subt) //插入一条子任务
	return anaType
}

// 检测是否有失败解析的子任务
func CheckFailedAnalyzeData() {
	for {
		time.Sleep(1 * time.Hour) //每隔一小时进行检查一次
		data := RabbitMqServer.GetData(failedquePath)
		if data != "" {
			splitdata := strings.Split(data, "?")[1]
			sendMsg := "requestanalyze?" + splitdata
			//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
			for _, val := range allAnalyzeClient {
				n, err := GetConn(val.Ip, "anaclient").Write([]byte(sendMsg))
				if err != nil && n == 0 {
					Logs.Loggers().Print("发送重新解析消息失败----", err.Error(), data)
					RabbitMqServer.PutData(failedquePath, data)
					break
				} else {
					//
					Logs.Loggers().Print("发送重新解析消息成功----", sendMsg)
					break
				}
			}
		}
	}
}
