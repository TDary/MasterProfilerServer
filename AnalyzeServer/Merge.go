package AnalyzeServer

import (
	"MasterServer/Data"
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
	"MasterServer/RabbitMqServer"
	"MasterServer/Tools"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

//对解析成功的消息进行检查判断是否可以进行合并操作
func AnalyzeSuccessToMerge() {
	for {
		time.Sleep(30 * time.Second) //每隔30秒检查一次
		CheckCaseToMerge()
	}
}

//合并simple数据
func MergeSimple(maintable DataBase.MainTable, dataPath string) {
	var allSimpleData Data.Simples
	var insertDatas []DataBase.InsertSimple
	for _, val := range maintable.RawFiles {
		rawPath := dataPath + "/" + val
		var isExit, _ = ioutil.ReadFile(rawPath)
		currentSimpleData := &Data.Simples{}
		if isExit != nil {
			//存在分析文件，可直接反序列化,解压后再反序列化
			err := Tools.ExtractZip(rawPath, dataPath)
			if err != nil {
				//解压失败
				Logs.Loggers().Print("解压分析文件失败----", rawPath)
				return
			}
			simpleDataPath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw.csv"
			bytedata, err := ioutil.ReadFile(simpleDataPath)
			if err != nil {
				//打开失败
				Logs.Loggers().Print("打开分析文件失败----", simpleDataPath)
				return
			}
			err = proto.Unmarshal(bytedata, currentSimpleData)
			if err != nil {
				Logs.Loggers().Print("反序列化失败----", err.Error())
				return
			}
			for key, val := range currentSimpleData.Allvalues {
				ishasdata := false
				for key1 := range allSimpleData.Allvalues {
					if key == key1 {
						ishasdata = true
						allSimpleData.Allvalues[key1].Values = append(allSimpleData.Allvalues[key1].Values, val.GetValues()...)
					}
				}
				if !ishasdata {
					allSimpleData.Allvalues = append(allSimpleData.Allvalues, currentSimpleData.Allvalues[key])
				}
			}
		} else {
			//不存在分析文件，先从minio下载在进行反序列化
			objectName := maintable.UUID + "/" + val
			isdownloadSuccess := Minio.DownLoadFile(objectName, rawPath, "application/zip")
			if isdownloadSuccess {
				err := Tools.ExtractZip(rawPath, dataPath)
				if err != nil {
					//解压失败
					Logs.Loggers().Print("解压分析文件失败----", rawPath)
					return
				}
				simpleDataPath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw.csv"
				bytedata, err := ioutil.ReadFile(simpleDataPath)
				if err != nil {
					//打开失败
					Logs.Loggers().Print("打开分析文件失败----", rawPath)
					return
				}
				err = proto.Unmarshal(bytedata, currentSimpleData)
				if err != nil {
					Logs.Loggers().Print("反序列化失败----", err.Error())
					return
				}
				for key, val := range currentSimpleData.Allvalues {
					ishasdata := false
					for key1, _ := range allSimpleData.Allvalues {
						if key == key1 {
							ishasdata = true
							allSimpleData.Allvalues[key1].Values = append(allSimpleData.Allvalues[key1].Values, val.GetValues()...)
						}
					}
					if !ishasdata {
						allSimpleData.Allvalues = append(allSimpleData.Allvalues, val)
					}
				}
			} else {
				Logs.Loggers().Printf("下载分析文件失败%s,UUID:%s----", objectName, maintable.UUID)
				return
			}
			//下载后解压进行反序列化
		}
	}
	//合并完成，转换数据结构进行入库操作
	for _, val := range allSimpleData.Allvalues {
		var item DataBase.InsertSimple
		item.Name = val.Field
		if item.Name == "frame" {
			item.Values = val.Values
			item.UUID = maintable.UUID
			var frame float32
			frame = 1
			for key, _ := range item.Values {
				item.Values[key] = frame //赋于正确的帧数
				frame += 1
			}
		} else {
			item.UUID = maintable.UUID
			item.Values = val.Values
		}
		insertDatas = append(insertDatas, item)
	}
	//入库
	DataBase.InsertSimpleData(insertDatas)
	DataBase.ModifyMain(maintable.UUID, 1, len(insertDatas[0].Values))
	//成功合并上报
	Tools.SendRobotMsg(config.RobotUrl, "UUID:"+maintable.UUID+"案例解析合并完成")
}

//合并funprofiler数据
func MergeFun(maintable DataBase.MainTable, dataPath string) {
	var allFunRow Data.AllCaseFunRow
	var allCaseFunName Data.ListCaseFunName
	var insertCaseFunRow []DataBase.CaseFunRow
	var insertCaseFunPath DataBase.CaseFunNamePath
	var frame int32
	frame = 1
	for _, val := range maintable.RawFiles {
		rawPath := dataPath + "/" + val
		currentFunRowData := &Data.AllCaseFunRow{}
		currentFunNameData := &Data.ListCaseFunName{}
		var isExit, _ = ioutil.ReadFile(rawPath)
		if isExit != nil {
			//存在分析文件，可直接反序列化,解压后再反序列化
			err := Tools.ExtractZip(rawPath, dataPath)
			if err != nil {
				//解压失败
				Logs.Loggers().Print("解压分析文件失败----", rawPath)
				return
			}
			rowDataPath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw_funrow.bin"
			funNamePath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw_funname.bin"
			//合并FunRow
			bytedata, err := ioutil.ReadFile(rowDataPath)
			if err != nil {
				//打开失败
				Logs.Loggers().Print("打开分析文件失败----", rowDataPath)
				return
			}
			err = proto.Unmarshal(bytedata, currentFunRowData)
			if err != nil {
				Logs.Loggers().Print("反序列化失败----", err.Error())
				return
			}
			var subCount int32
			for _, val := range currentFunRowData.Allvalues {
				if val.Name == "Main Thread" {
					subCount = 300 - int32(len(val.Frames)) //默认300帧中会少2帧，但是由于某些情况中间有可能会缺一些
				}
				ishasdata := false
				for key2, val2 := range allFunRow.Allvalues {
					if val2.Name == val.Name {
						ishasdata = true
						for ks, fs := range val.Frames {
							val.Frames[ks].Frame = frame + fs.Frame - subCount
						}
						allFunRow.Allvalues[key2].Frames = append(allFunRow.Allvalues[key2].Frames, val.Frames...)
					}
				}
				if !ishasdata {
					for key, fs := range val.Frames {
						val.Frames[key].Frame = frame + fs.Frame - subCount
					}
					allFunRow.Allvalues = append(allFunRow.Allvalues, val)
				}
			}
			//合并FunNamePath
			bytedata2, err := ioutil.ReadFile(funNamePath)
			if err != nil {
				//打开失败
				Logs.Loggers().Print("打开分析文件失败----", funNamePath)
				return
			}
			err = proto.Unmarshal(bytedata2, currentFunNameData)
			if err != nil {
				Logs.Loggers().Print("反序列化失败----", err.Error())
				return
			}
			for _, name := range currentFunNameData.Funnames {
				ishasdata := false
				for _, name2 := range allCaseFunName.Funnames {
					if name == name2 {
						ishasdata = true
						//不重复添加，在此跳过
					}
				}
				if !ishasdata {
					allCaseFunName.Funnames = append(allCaseFunName.Funnames, name)
				}
			}
		} else {
			//不存在分析文件，先从minio下载在进行反序列化
			objectName := maintable.UUID + "/" + val
			isdownloadSuccess := Minio.DownLoadFile(objectName, rawPath, "application/zip")
			if isdownloadSuccess {
				err := Tools.ExtractZip(rawPath, dataPath)
				if err != nil {
					//解压失败
					Logs.Loggers().Print("解压分析文件失败----", rawPath)
					return
				}
				rowDataPath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw_funrow.bin"
				funNamePath := dataPath + "/" + strings.Split(val, ".")[0] + ".raw_funname.bin"
				//合并FunRow
				bytedata, err := ioutil.ReadFile(rowDataPath)
				if err != nil {
					//打开失败
					Logs.Loggers().Print("打开分析文件失败----", rowDataPath)
					return
				}
				err = proto.Unmarshal(bytedata, currentFunRowData)
				if err != nil {
					Logs.Loggers().Print("反序列化失败----", err.Error())
					return
				}
				var subCount int32
				for _, val := range currentFunRowData.Allvalues {
					if val.Name == "Main Thread" {
						subCount = 300 - int32(len(val.Frames)) //默认300帧中会少2帧，但是由于某些情况中间有可能会缺一些
					}
					ishasdata := false
					for key2, val2 := range allFunRow.Allvalues {
						if val2.Name == val.Name {
							ishasdata = true
							for ks, fs := range val.Frames {
								val.Frames[ks].Frame = frame + fs.Frame - subCount
							}
							allFunRow.Allvalues[key2].Frames = append(allFunRow.Allvalues[key2].Frames, val.Frames...)
						}
					}
					if !ishasdata {
						for key, fs := range val.Frames {
							val.Frames[key].Frame = frame + fs.Frame - subCount
						}
						allFunRow.Allvalues = append(allFunRow.Allvalues, val)
					}
				}
				//合并FunNamePath
				bytedata2, err := ioutil.ReadFile(funNamePath)
				if err != nil {
					//打开失败
					Logs.Loggers().Print("打开分析文件失败----", funNamePath)
					return
				}
				err = proto.Unmarshal(bytedata2, currentFunNameData)
				if err != nil {
					Logs.Loggers().Print("反序列化失败----", err.Error())
					return
				}
				for _, name := range currentFunNameData.Funnames {
					ishasdata := false
					for _, name2 := range allCaseFunName.Funnames {
						if name == name2 {
							ishasdata = true
							//不重复添加，在此跳过
						}
					}
					if !ishasdata {
						allCaseFunName.Funnames = append(allCaseFunName.Funnames, name)
					}
				}
			}
		}
		if len(currentFunRowData.Allvalues[0].Frames) != 0 {
			count := len(currentFunRowData.Allvalues[0].Frames)
			frame += int32(count)
		}
	}
	//转换数据结构
	var totalFrame int
	for _, vals := range allFunRow.Allvalues {
		var caseFunRow DataBase.CaseFunRow
		caseFunRow.UUID = maintable.UUID
		caseFunRow.Name = vals.Name
		if caseFunRow.Name == "Main Thread" { //暂时不放入这个数据，以后要的话再说
			totalFrame = len(vals.Frames)
			continue
		}
		var totalTime int32
		for _, va2 := range vals.Frames {
			var funrowInfo DataBase.FunRowInfo
			funrowInfo.Frame = va2.Frame
			funrowInfo.Total = va2.Total
			funrowInfo.Self = va2.Self
			funrowInfo.Calls = va2.Calls
			funrowInfo.Gcalloc = va2.Gcalloc
			funrowInfo.Timems = va2.Timems
			funrowInfo.Selfms = va2.Selfms
			totalTime += va2.Timems
			caseFunRow.Frames = append(caseFunRow.Frames, funrowInfo)
		}
		caseFunRow.AvgValidTime = totalTime / int32(len(vals.Frames))
		insertCaseFunRow = append(insertCaseFunRow, caseFunRow)
	}
	insertCaseFunPath.UUID = maintable.UUID
	insertCaseFunPath.Stack = allCaseFunName.Funnames
	//入库
	DataBase.InsertCaseFunRow(insertCaseFunRow)
	DataBase.InsertFunNamePath(insertCaseFunPath)
	DataBase.ModifyMain(maintable.UUID, 1, totalFrame)
	//成功合入库上报
	Tools.SendRobotMsg(config.RobotUrl, "UUID:"+maintable.UUID+"案例解析合并完成")
}

//开始合并且入库操作
func MergeBegin(maintable DataBase.MainTable) {
	dataPath := config.MergePath + "/" + maintable.UUID
	_, err := os.Stat(dataPath)
	if err != nil {
		Logs.Loggers().Printf("当前文件夹%s不存在，重新创建中！", dataPath)
		os.Mkdir(dataPath, 0755)
	}
	if maintable.AnalyzeType == "simple" {
		//simple数据合并
		MergeSimple(maintable, dataPath)
	} else if maintable.AnalyzeType == "funprofiler" {
		//funprofiler合并
		MergeFun(maintable, dataPath)
	} else {
		//deep合并
	}
}

//检查案例状态是否有可以进行合并的
func CheckCaseToMerge() {
	waitCase := DataBase.FindMainTable(0)
	if len(waitCase) > 0 {
		for _, val := range waitCase {
			currentCase := CheckSub(val.UUID)
			if currentCase {
				//当前案例可以进行合并操作了
				DataBase.ModifyMainState(val.UUID, 3)
				Logs.Loggers().Print("开始合并案例,UUID:" + val.UUID)
				go MergeBegin(val)
			} else {
				continue
			}
		}
	}
}

//检查子表
func CheckSub(uuid string) bool {
	subt := DataBase.FindSubTableData(uuid)
	if subt != nil {
		for _, val := range subt {
			if val.State == 1 {
				continue
			} else {
				return false
			}
		}
		return true
	} else {
		Logs.Loggers().Print("不存在该UUID：" + uuid + "的子表数据")
		return false
	}
}

//处理解析成功消息状态
func ParseSuccessData(data string) {
	var addData SuccessData
	splidata := strings.Split(data, "&")
	for i := 0; i < len(splidata); i++ {
		if strings.Contains(splidata[i], "ip") {
			current := strings.Split(splidata[i], "=")
			addData.IP = current[1]
		} else if strings.Contains(splidata[i], "rawfile") {
			current := strings.Split(splidata[i], "=")
			addData.RawFile = current[1]
		} else if strings.Contains(splidata[i], "uuid") {
			current := strings.Split(splidata[i], "=")
			addData.UUID = current[1]
		}
	}
	if strings.Contains(addData.RawFile, ".zip") { //转换一下源文件
		addData.RawFile = strings.Split(addData.RawFile, ".")[0] + ".raw"
	}
	DataBase.UpdateStates(addData.RawFile, addData.UUID, 1, addData.IP) //更新状态值
}

//处理解析失败消息状态
func ParseFailedData(data string) {
	var addData SuccessData
	splidata := strings.Split(data, "&")
	for i := 0; i < len(splidata); i++ {
		if strings.Contains(splidata[i], "ip") {
			current := strings.Split(splidata[i], "=")
			addData.IP = current[1]
		} else if strings.Contains(splidata[i], "rawfile") {
			current := strings.Split(splidata[i], "=")
			addData.RawFile = current[1]
		} else if strings.Contains(splidata[i], "uuid") {
			current := strings.Split(splidata[i], "=")
			addData.UUID = current[1]
		}
	}
	DataBase.UpdateStates(addData.RawFile, addData.UUID, -1, addData.IP) //更新状态值
	failedquePath := "./ServerQue/" + "FialedAnalyzeQue"
	RabbitMqServer.PutData(failedquePath, data)
}
