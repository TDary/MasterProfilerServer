package AnalyzeServer

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"MasterServer/Minio"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//停止采集信号
func StopAnalyzeRequest(mes string) {
	var data EndData
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		if strings.Contains(str1[i], "uuid") { //解析gameid
			uid := strings.Split(str1[i], "=")
			data.UUID = uid[1]
		} else if strings.Contains(str1[i], "ip") { //解析uuid
			currentIp := strings.Split(str1[i], "=")
			data.Ip = currentIp[1]
		} else if strings.Contains(str1[i], "lastfile") { //解析rawfiles
			file := strings.Split(str1[i], "=")
			data.LastRawFile = file[1]
		} else {
			Logs.Loggers().Print("不存在对额外参数的解析:" + str1[i])
			Logs.Loggers().Print("无法识别的完整参数：" + mes)
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

//解析url进行结构化并创建数据库表数据
func AnalyzeRequest(data string) {
	//此处作为消费者,同时调用DataBase创建数据库表
	getdata, mtable := ReceiveMes(data)
	var rawFiles []string
	for {
		//检测是否有停止解析请求,开一个数组，保存停止解析消息
		state, file := GetLastRawFileIsSend(getdata.UUID, rawFiles)
		if state == 0 {
			//所有任务都已完成，在此打断循环断开采集任务
			//更新主表数据库
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
				//更新数据库主表
				DataBase.UpdateMainTable(getdata.Appkey, getdata.UUID, rawFiles)
				break
			}
			break //打断循环
		} else {
			//wait
		}
		time.Sleep(10 * time.Second)                       //每隔10秒检测一次
		rFiles := Minio.SearchObjectOfBucket(getdata.UUID) //检测 并发送解析请求
		for i := 0; i < len(rFiles); i++ {
			getdata.RawFileName = rFiles[i]
			getdata.RawFile = strings.Split(rFiles[i], "/")[1]
			if len(rawFiles) == 0 {
				rawFiles = append(rawFiles, getdata.RawFile)
				//发送解析请求,随便发送一台空闲的解析器让其进行轮转解析
				GetIdleAnalyzeClient()
				for _, val := range allAnalyzeClient {
					SendRequestAnalyze(getdata, val.IpAddress)
					//插入数据库子任务
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
						InsertSubTable(mtable, getdata.RawFile)
						break
					}
				}
			}
		}
	}
}

//对解析成功的消息进行检查判断是否可以进行合并操作
func AnalyzeSuccessUrl() {
	//也是进行轮询查找,一次查找较多的数据
	var getanalyzeData string
	var waitModifyState []SuccessData
	isAnalyzeStop = true //初始状态系统关闭(避免刷日志)
	for {
		if !isAnalyzeStop {
			if len(waitModifyState) == 50 {
				//达到了允许存储的上限,直接进行修改状态值并释放进程
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

//检查案例状态是否有可以进行合并的
func CheckCaseToMerge() {
	var waitCase []DataBase.MainTable
	waitCase = DataBase.FindMainTable(0)
	if waitCase != nil {
		for i := 0; i < len(waitCase); i++ {
			currentCase := CheckSub(waitCase[i].UUID)
			if currentCase {
				//当前案例可以进行合并操作了
				//发送合并消息给合并服务器
				SendMergeData(waitCase[i].UUID)
			} else {
				continue
			}
		}
	} else {
		Logs.Loggers().Print("无待合并的案例----")
	}
}

//发送请求合并消息
func SendMergeData(uuid string) {
	request_Url := "http://" + config.MergeServer.Ip + ":" + config.MergeServer.Port +
		"/merge" + "?" + "uuid=" + uuid
	//超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(request_Url)
	if err != nil {
		Logs.Loggers().Print(err)
		return
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			Logs.Loggers().Print(err)
		}
	}
	if result.String() == "ok" {
		Logs.Loggers().Print("UUID：" + uuid + "接收到合并消息，即将开始合并入库操作----")
	} else {
		Logs.Loggers().Print("客户端未成功接收到消息----")
		return
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

//添加客户端解析器进入组网
func AddAnalyzeClient(data string) {
	//当客户端解析器启动时会ping一次服务器，测试是否已将客户端解析器加入了组网,保证机器间正常运行
	//只有当ping通的情况下才会将开关打开
	//同时可以加入新的解析客户端
	//同时在此处进行启动客户端解析
	nowstr := strings.Split(data, "&")
	nowip := strings.Split(nowstr[0], "=")
	if CheckKey(nowip[1]) {
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
