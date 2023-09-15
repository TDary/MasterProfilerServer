package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
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
func GetLastRawFileIsSend(uuid string, rawFiles []string) (int, string) {
	for _, val := range stopMsg {
		if val.UUID == uuid {
			//在停止消息中有这份报告，再检测最后一份源文件是否已发送解析
			for _, vals := range rawFiles {
				if vals == val.LastRawFile {
					//已经有了且已发送解析，此时不再进行下去
					return 0, ""
				}
			}
			//没有发送，还要再发最后一次
			return 1, val.LastRawFile
		}
	}
	return -1, "" //继续等待
}

//将原始文件进行排序
func SortRawFils(rawfiles []string) {
	for i := 0; i < len(rawfiles)-1; i++ {
		for j := 0; j < len(rawfiles)-i-1; j++ {
			if strings.Split(rawfiles[j], ".")[0] > strings.Split(rawfiles[j+1], ".")[0] {
				rawfiles[j], rawfiles[j+1] = rawfiles[j+1], rawfiles[j]
			}
		}
	}
}

//解析url进行结构化并创建数据库表数据
func AnalyzeRequest(data string) {
	//此处作为消费者,同时调用DataBase创建数据库表
	getdata, mtable := ReceiveMes(data)
	var rawFiles []string
	quePath := "./ServerQue/" + getdata.UUID + "_AnalyzeQue"
	for {
		analze := GetSuccessMes(quePath)
		if analze != "" {
			res := GetAnalyzeData(analze)
			//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
			GetIdleAnalyzeClient()
			for _, val := range allAnalyzeClient {
				SendRequestAnalyze(res, val.Ip)
				//插入数据库子任务
				InsertSubTable(mtable, res.RawFile)
				//更新数据库主表，先排序，这里用了冒泡
				//SortRawFils(rawFiles)
				break
			}
			rawFiles = append(rawFiles, res.RawFile)
		}
		state, file := GetLastRawFileIsSend(getdata.UUID, rawFiles)
		if state == 0 {
			//所有任务都已完成，在此打断循环断开采集任务
			//更新主表数据库
			//SortRawFils(rawFiles) //排一下序，因为原有的插入顺序不一定正确
			DataBase.UpdateMainTable(getdata.Appkey, getdata.UUID, rawFiles)
			break
		} else if state == 1 && file != "" {
			//还有最后一个没发送解析,在此发送
			rawFiles = append(rawFiles, file)
			//队列中拿rawfilename
			getanalyzeData := GetSuccessMes(quePath)
			if getanalyzeData != "" {
				res := GetAnalyzeData(getanalyzeData)
				if res.RawFile == file {
					getdata.RawFile = res.RawFile
					getdata.RawFileName = res.RawFileName
					//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
					GetIdleAnalyzeClient()
					for _, val := range allAnalyzeClient {
						SendRequestAnalyze(getdata, val.Ip)
						//插入数据库子任务
						InsertSubTable(mtable, getdata.RawFile)
						//更新数据库主表，先排序，这里用了冒泡
						//SortRawFils(rawFiles)
						DataBase.UpdateMainTable(getdata.Appkey, getdata.UUID, rawFiles)
						break
					}
					break
				}
			}
		}
	}
	for {
		rFiles := Minio.SearchObjectOfBucket(getdata.UUID) //检测 并发送解析请求
		for i := 0; i < len(rFiles); i++ {
			Logs.Loggers().Print("rawfileName:" + rFiles[i])
			getdata.RawFileName = rFiles[i]                    //uuid/1231.zip
			getdata.RawFile = strings.Split(rFiles[i], "/")[1] //21313.zip
			if len(rawFiles) == 0 {
				rawFiles = append(rawFiles, getdata.RawFile)
				//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
				GetIdleAnalyzeClient()
				for _, val := range allAnalyzeClient {
					SendRequestAnalyze(getdata, val.IpAddress)
					//插入数据库子任务
					isMergeStop = false
					InsertSubTable(mtable, getdata.RawFile)
					break
				}
			} else {
				//检测是否含有已发送的任务
				isHasSame := false
				for i := 0; i < len(rawFiles); i++ {
					if rawFiles[i] == getdata.RawFile {
						isHasSame = true
						break
					}
				}
				if !isHasSame {
					//该源没有相同的，可进行发送解析
					GetIdleAnalyzeClient()
					for _, val := range allAnalyzeClient {
						SendRequestAnalyze(getdata, val.IpAddress)
						//插入数据库子任务
						isMergeStop = false
						InsertSubTable(mtable, getdata.RawFile)
						break
					}
				}
			}
		}
		//检测是否有停止解析请求,开一个数组，保存停止解析消息
		state, file := GetLastRawFileIsSend(getdata.UUID, rawFiles)
		if state == 0 {
			//所有任务都已完成，在此打断循环断开采集任务
			//更新主表数据库
			SortRawFils(rawFiles) //排一下序，因为原有的插入顺序不一定正确
			DataBase.UpdateMainTable(getdata.Appkey, getdata.UUID, rawFiles)
			break
		} else if state == 1 && file != "" {
			//还有最后一个没发送解析,在此发送
			rawFiles = append(rawFiles, file)
			getdata.RawFile = file
			rFiles := Minio.SearchObjectOfBucket(getdata.UUID)
			for _, val := range rFiles {
				if strings.Split(val, "/")[1] == file {
					getdata.RawFileName = val
					break
				}
			}
			//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
			GetIdleAnalyzeClient()
			for _, val := range allAnalyzeClient {
				SendRequestAnalyze(getdata, val.IpAddress)
				//插入数据库子任务
				InsertSubTable(mtable, getdata.RawFile)
				//更新数据库主表，先排序，这里用了冒泡
				SortRawFils(rawFiles)
				DataBase.UpdateMainTable(getdata.Appkey, getdata.UUID, rawFiles)
				break
			}
			break //打断循环
		} else {
			//wait
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
