package AnalyzeServer

import (
	"MasterServer/Data"
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
	"MasterServer/Tools"
	"io/ioutil"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

//对解析成功的消息进行检查判断是否可以进行合并操作
func AnalyzeSuccessUrl() {
	//也是进行轮询查找,一次查找较多的数据
	var getanalyzeData string
	var waitModifyState []SuccessData
	for {
		if !isMergeStop {
			if len(waitModifyState) == 50 {
				//达到了允许存储的上限,直接进行修改状态值
				var allip []string
				ModifySubState(waitModifyState, allip) //修改状态值
				waitModifyState = nil                  //重置上限值
				//开始判断是否有案例可以进行合并入库操作
				CheckCaseToMerge()
			}
			getanalyzeData = GetSuccessMes("./ServerQue/ParseQueSuccessQue")
			if getanalyzeData != "" {
				waitModifyState = ParseSuccessData(getanalyzeData, waitModifyState)
			} else {
				Logs.Loggers().Print("成功解析消息队列已空，进入检查状态")
				//开始修改子案例状态,同时释放解析进程
				isMergeStop = true
				var allip []string
				ModifySubState(waitModifyState, allip) //修改状态值
				waitModifyState = nil                  //重置上限值
				//开始判断是否有案例可以进行合并入库操作
				CheckCaseToMerge()
			}
		}
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
			raw, err := Tools.ExtractZip(rawPath, dataPath)
			if err != nil {
				//解压失败
				Logs.Loggers().Print("解压分析文件失败----", rawPath)
				return
			}
			rawDataPath := dataPath + "/" + raw
			bytedata, err := ioutil.ReadFile(rawDataPath)
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
				raw, err := Tools.ExtractZip(rawPath, dataPath)
				if err != nil {
					//解压失败
					Logs.Loggers().Print("解压分析文件失败----", rawPath)
					return
				}
				rawDataPath := dataPath + "/" + raw
				bytedata, err := ioutil.ReadFile(rawDataPath)
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
			item.Valus = val.Values
			item.UUID = maintable.UUID
			var frame float32
			frame = 1
			for key, _ := range item.Valus {
				item.Valus[key] = frame //赋于正确的帧数
				frame += 1
			}
		} else {
			item.Valus = val.Values
		}
		insertDatas = append(insertDatas, item)
	}
	//入库
	DataBase.InsertSimpleData(insertDatas)
	DataBase.ModifyMainState(maintable.AppKey, maintable.UUID, 1)
}

//合并funprofiler数据
func MergeFun(maintable DataBase.MainTable, dataPath string) {
	var allFunRow Data.AllCaseFunRow
	var insertCaseFunRow []DataBase.CaseFunRow
	for _, val := range maintable.RawFiles {
		//objectName = val
		rawPath := dataPath + "/" + val
		var isExit, _ = ioutil.ReadFile(rawPath)
		currentFunRowData := &Data.AllCaseFunRow{}
		if isExit != nil {
			//存在分析文件，可直接反序列化
			//存在分析文件，可直接反序列化,解压后再反序列化
			raw, err := Tools.ExtractZip(rawPath, dataPath)
			if err != nil {
				//解压失败
				Logs.Loggers().Print("解压分析文件失败----", rawPath)
				return
			}
			rawDataPath := dataPath + "/" + raw
			bytedata, err := ioutil.ReadFile(rawDataPath)
			if err != nil {
				//打开失败
				Logs.Loggers().Print("打开分析文件失败----", rawPath)
				return
			}
			err = proto.Unmarshal(bytedata, currentFunRowData)
			if err != nil {
				Logs.Loggers().Print("反序列化失败----", err.Error())
				return
			}
			for _, val := range currentFunRowData.Allvalues {
				ishasdata := false
				for key2, val2 := range allFunRow.Allvalues {
					if val2.Name == val.Name {
						ishasdata = true
						allFunRow.Allvalues[key2].Frames = append(allFunRow.Allvalues[key2].Frames, val.Frames...)
					}
				}
				if !ishasdata {
					allFunRow.Allvalues = append(allFunRow.Allvalues, val)
				}
			}
		} else {
			//不存在分析文件，先从minio下载在进行反序列化
			objectName := maintable.UUID + "/" + val
			isdownloadSuccess := Minio.DownLoadFile(objectName, rawPath, "application/zip")
			if isdownloadSuccess {
				raw, err := Tools.ExtractZip(rawPath, dataPath)
				if err != nil {
					//解压失败
					Logs.Loggers().Print("解压分析文件失败----", rawPath)
					return
				}
				rawDataPath := dataPath + "/" + raw
				bytedata, err := ioutil.ReadFile(rawDataPath)
				if err != nil {
					//打开失败
					Logs.Loggers().Print("打开分析文件失败----", rawPath)
					return
				}
				err = proto.Unmarshal(bytedata, currentFunRowData)
				if err != nil {
					Logs.Loggers().Print("反序列化失败----", err.Error())
					return
				}
				for _, val := range currentFunRowData.Allvalues {
					ishasdata := false
					for key2, val2 := range allFunRow.Allvalues {
						if val2.Name == val.Name {
							ishasdata = true
							allFunRow.Allvalues[key2].Frames = append(allFunRow.Allvalues[key2].Frames, val.Frames...)
						}
					}
					if !ishasdata {
						allFunRow.Allvalues = append(allFunRow.Allvalues, val)
					}
				}
			}
		}
	}
	//转换数据结构
	for _, vals := range allFunRow.Allvalues {
		var caseFunRow DataBase.CaseFunRow
		frame := 1
		caseFunRow.UUID = maintable.UUID
		caseFunRow.Name = vals.Name
		if caseFunRow.Name == "Main Thread" { //暂时不要放入这个数据，以后要的话再说
			continue
		}
		for _, va2 := range vals.Frames {
			var funrowInfo DataBase.FunRowInfo
			funrowInfo.Frame = int32(frame)
			funrowInfo.Total = va2.Total
			funrowInfo.Self = va2.Self
			funrowInfo.Calls = va2.Calls
			funrowInfo.Gcalloc = va2.Gcalloc
			funrowInfo.Timems = va2.Timems
			funrowInfo.Selfms = va2.Selfms
			caseFunRow.Frames = append(caseFunRow.Frames, funrowInfo)
			frame += 1
		}
		insertCaseFunRow = append(insertCaseFunRow, caseFunRow)
	}
	//入库
	DataBase.InsertCaseFunRow(insertCaseFunRow)
	DataBase.ModifyMainState(maintable.AppKey, maintable.UUID, 1)
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
	var waitCase []DataBase.MainTable
	waitCase = DataBase.FindMainTable(0)
	if waitCase != nil {
		for _, val := range waitCase {
			currentCase := CheckSub(val.UUID)
			if currentCase {
				//当前案例可以进行合并操作了
				DataBase.ModifyMainState(val.AppKey, val.UUID, 3)
				Logs.Loggers().Print("开始合并案例,UUID:" + val.UUID)
				go MergeBegin(val)
			} else {
				continue
			}
		}
	} else {
		Logs.Loggers().Print("无待合并的案例----")
	}
	waitCase = nil //上面流程完毕清除一次
}

//检查子表
func CheckSub(uuid string) bool {
	subt := DataBase.FindSubTableData(uuid)
	if subt != nil {
		for i := 0; i < len(subt); i++ {
			if subt[i].State == 0 {
				return false
			}
		}
		return true
	} else {
		Logs.Loggers().Print("不存在该UUID：" + uuid + "的子表数据")
		return false
	}
}

//修改子案例状态
func ModifySubState(wdata []SuccessData, allip []string) {
	for _, val := range wdata {
		allip = append(allip, val.IP)
		DataBase.UpdateStates(val.RawFile, val.UUID, 1, val.IP)
	}
}

//处理解析成功消息
func ParseSuccessData(data string, wdata []SuccessData) []SuccessData {
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
	wdata = append(wdata, addData)
	return wdata
}
