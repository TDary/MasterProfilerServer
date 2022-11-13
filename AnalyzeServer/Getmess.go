package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"strings"
)

func ReceiveMes(mes string) {
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
			fs := strings.Split(files[1], ",")
			mtable.RawFiles = fs
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
			Logs.Loggers().Print("无法识别的完整参数：" + mes)
			return
		}
	}
	mtable.State = 0
	mtable.ScreenFiles = nil
	mtable.ScreenState = 0
	DataBase.InsertMain(mtable)
	GetSubData(mtable)
}

func GetSubData(mtable DataBase.MainTable) {
	if len(mtable.RawFiles) == 0 {
		return
	}
	for i := 0; i < len(mtable.RawFiles); i++ {
		var stable DataBase.SubTable
		stable.AppKey = mtable.AppKey
		stable.UUID = mtable.UUID
		stable.State = mtable.State
		stable.Priority = mtable.Priority
		stable.StorageIp = mtable.StorageIp
		stable.UnityVersion = mtable.UnityVersion
		stable.AnalyzeBucket = mtable.AnalyzeBucket
		stable.AnalyzeIP = ""
		stable.CsvPath = ""
		stable.RawFile = mtable.RawFiles[i]
		DataBase.InsertSub(stable)
	}
}
