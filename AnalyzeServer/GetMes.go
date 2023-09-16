package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/RabbitMqServer"
	"strings"
)

//从队列中拿出数据
func GetSuccessMes(data string) string {
	res := RabbitMqServer.GetData(data)
	return res
}

func GetAnalyzeData(data string) AnalyzeData {
	var res AnalyzeData
	str1 := strings.Split(data, "&")
	for i := 0; i < len(str1); i++ {
		if strings.Contains(str1[i], "uuid") { //解析uuid
			uid := strings.Split(str1[i], "=")
			res.UUID = uid[1]
		} else if strings.Contains(str1[i], "rawfile") { //解析rawfiles
			file := strings.Split(str1[i], "=")
			res.RawFile = file[1]
		} else if strings.Contains(str1[i], "rawfilename") { //解析gamename
			na := strings.Split(str1[i], "=")
			res.RawFileName = na[1]
		} else if strings.Contains(str1[i], "unityVersion") { //解析unityVersion
			na := strings.Split(str1[i], "=")
			res.UnityVersion = na[1]
		} else if strings.Contains(str1[i], "analyzebucket") { //解析AnalyzeBucket
			na := strings.Split(str1[i], "=")
			res.Bucket = na[1]
		} else if strings.Contains(str1[i], "analyzeType") { //解析类型
			na := strings.Split(str1[i], "=")
			res.AnalyzeType = na[1]
		}
	}
	return res
}

//接受开始采集消息
func ReceiveMes(mes string) DataBase.MainTable {
	var mtable DataBase.MainTable
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		if strings.Contains(str1[i], "gameid") { //解析gameid
			gid := strings.Split(str1[i], "=")
			mtable.AppKey = gid[1]
		} else if strings.Contains(str1[i], "uuid") { //解析uuid
			uid := strings.Split(str1[i], "=")
			mtable.UUID = uid[1]
		} else if strings.Contains(str1[i], "rawFiles") { //解析rawfiles
			files := strings.Split(str1[i], "=")
			if files[1] != "" {
				fs := strings.Split(files[1], ",")
				mtable.RawFiles = fs
			}
		} else if strings.Contains(str1[i], "gameName") { //解析gamename
			na := strings.Split(str1[i], "=")
			mtable.GameName = na[1]
		} else if strings.Contains(str1[i], "caseName") { //解析casename
			na := strings.Split(str1[i], "=")
			mtable.CaseName = na[1]
		} else if strings.Contains(str1[i], "unityVersion") { //解析unityVersion
			na := strings.Split(str1[i], "=")
			mtable.UnityVersion = na[1]
		} else if strings.Contains(str1[i], "bucket") { //解析AnalyzeBucket
			na := strings.Split(str1[i], "=")
			mtable.AnalyzeBucket = na[1]
		} else if strings.Contains(str1[i], "anatype") { //解析类型
			na := strings.Split(str1[i], "=")
			mtable.AnalyzeType = na[1]
		} else if strings.Contains(str1[i], "storageIp") { //解析StorageIp
			na := strings.Split(str1[i], "=")
			mtable.StorageIp = na[1]
		} else if strings.Contains(str1[i], "device") { //解析Device
			na := strings.Split(str1[i], "=")
			mtable.Device = na[1]
		} else if strings.Contains(str1[i], "beginTime") { //解析TestBeginTime
			na := strings.Split(str1[i], "=")
			mtable.TestBeginTime = na[1]
		} else if strings.Contains(str1[i], "endTime") { //解析TestEndTime
			na := strings.Split(str1[i], "=")
			mtable.TestEndTime = na[1]
		} else if strings.Contains(str1[i], "priority") { //解析优先级priority
			na := strings.Split(str1[i], "=")
			mtable.Priority = na[1]
		} else {
			Logs.Loggers().Print("不存在对额外参数的解析:" + str1[i])
		}
	}
	mtable.State = 0
	mtable.ScreenFiles = nil
	mtable.ScreenState = 0
	if len(mtable.RawFiles) == 0 {
		mtable.RawFiles = nil
	}
	DataBase.InsertMain(mtable) //todo:判断是否有已经存在的Uuid，有的话不插入
	GetSubData(mtable)
	return mtable
}

func GetSubData(mtable DataBase.MainTable) {
	if len(mtable.RawFiles) == 0 {
		return
	}
	for i := 0; i < len(mtable.RawFiles); i++ {
		var stable DataBase.SubTable
		stable.UUID = mtable.UUID
		stable.State = mtable.State
		stable.AnalyzeIP = ""
		stable.RawFile = mtable.RawFiles[i]
		DataBase.InsertSub(stable)
	}
}

//插入一条子表任务
func InsertSubTable(mtable DataBase.MainTable, rawfile string) {
	var stable DataBase.SubTable
	stable.UUID = mtable.UUID
	stable.State = mtable.State
	stable.AnalyzeIP = ""
	stable.RawFile = rawfile
	DataBase.InsertSub(stable)
}

//插入一条子表任务
func InsertSubTableBySub(mtable DataBase.SubTable) {
	DataBase.InsertSub(mtable)
}
