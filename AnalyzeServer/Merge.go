package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"strings"
)

//对解析成功的消息进行检查判断是否可以进行合并操作
func AnalyzeSuccessUrl() {
	//也是进行轮询查找,一次查找较多的数据
	var getanalyzeData string
	var waitModifyState []SuccessData
	isAnalyzeStop = true //初始状态系统关闭(避免刷日志)
	for {
		if !isAnalyzeStop {
			if len(waitModifyState) == 50 {
				//达到了允许存储的上限,直接进行修改状态值
				var allip []string
				ModifySubState(waitModifyState, allip) //修改状态值
				waitModifyState = nil                  //重置上限值
				//开始判断是否有案例可以进行合并入库操作
				CheckCaseToMerge()
			}
			getanalyzeData = GetSuccessMes("/HttpServer/ParseQueSuccessQue")
			if getanalyzeData != "" {
				waitModifyState = ParseSuccessData(getanalyzeData, waitModifyState)
			} else {
				Logs.Loggers().Print("成功解析消息队列已空，进入检查状态")
				isAnalyzeStop = true
				//开始修改子案例状态,同时释放解析进程
				var allip []string
				ModifySubState(waitModifyState, allip) //修改状态值
				waitModifyState = nil                  //重置上限值
				//开始判断是否有案例可以进行合并入库操作
				CheckCaseToMerge()
			}
		}
	}
}

//开始合并且入库操作
func MergeBegin(maintable DataBase.MainTable) {

}

//检查案例状态是否有可以进行合并的
func CheckCaseToMerge() {
	var waitCase []DataBase.MainTable
	waitCase = DataBase.FindMainTable(0)
	if waitCase != nil {
		for i := 0; i < len(waitCase); i++ {
			currentCase := CheckSub(waitCase[i].UUID)
			if currentCase {
				//当前案例可以进行合并操作了
				Logs.Loggers().Print("开始合并案例,UUID:" + waitCase[i].UUID)
				MergeBegin(waitCase[i])
			} else {
				continue
			}
		}
	} else {
		Logs.Loggers().Print("无待合并的案例----")
	}
}

//检查子表
func CheckSub(uuid string) bool {
	var subt []DataBase.SubTable
	subt = DataBase.FindSubTableData(uuid)
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
