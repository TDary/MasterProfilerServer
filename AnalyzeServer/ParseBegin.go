package AnalyzeServer

import (
	"UAutoServer/DataBase"
	"UAutoServer/Logs"
	"bytes"
	"io"
	"net/http"
	"time"
)

func ParseEntrance() {
	for true {
		CheckSubTable()
		//每隔五分钟检查并发送一次解析
		//TODO:如果有请求案例进来则打断等待
		time.Sleep(5 * time.Minute)
	}
}

func SendBeginMessage(st DataBase.SubTable, m_ip string) {
	//发送开始解析的相关数据信息
	//前提是已经创建好数据表
	request_Url := "http://" + m_ip + "/analyze?uuid=" + st.UUID +
		"&rawfile=" + st.RawFile + "&unityversion=" + st.UnityVersion + "&analyzebucket=" + st.AnalyzeBucket
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
		//客户端成功接收到开始解析的消息，降空闲进程数减1
		ReduceRunC(m_ip, 1)
	} else {
		Logs.Loggers().Print("客户端未成功接收到消息----")
		return
	}
}

//查询获取到子表解析数量进行比对
func CheckSubTable() {
	var freeProcess int
	freeProcess = GetTotalRunC()
	currentState := 0
	highdata := DataBase.FindSTbyState(currentState) //高优先级的先解析
	andata := DataBase.FindSTbyState(currentState)   //普通案例
	if len(highdata) != 0 {
		Parse(freeProcess, highdata)
	} else {
		Parse(freeProcess, andata)
	}

}

//解析处理以及发送开始解析消息
func Parse(freec int, st []DataBase.SubTable) {
	if freec < len(st) { //可用进程数量小于子表解析数量
		for i := 0; i < freec; i++ {
			ip := GetFreeRunC()
			SendBeginMessage(st[i], ip)
		}
	} else if freec > len(st) {
		for i := 0; i < len(st); i++ {
			ip := GetFreeRunC()
			SendBeginMessage(st[i], ip)
		}
	} else {
		for i := 0; i < freec; i++ {
			ip := GetFreeRunC()
			SendBeginMessage(st[i], ip)
		}
	}
}

//每发送一次解析开始的请求就调用一次
func ReduceRunC(ip string, rcount int) {
	if CheckKey(ip) {
		if allclients[ip].WorkerNumbers == 0 {
			allclients[ip].State = false
			Logs.Loggers().Print("IP:" + ip + "的可用机器已用完，正在等待释放，该机器强行设置进入关闭状态。")
		} else {
			allclients[ip].WorkerNumbers = allclients[ip].WorkerNumbers - rcount
		}
	}
}

//获取总可用的解析进程
func GetTotalRunC() int {
	var totalCount int
	for _, value := range allclients {
		totalCount += value.WorkerNumbers
	}
	return totalCount
}

func GetFreeRunC() string {
	var freeClient string
	for _, val := range allclients {
		if val.WorkerNumbers != 0 {
			freeClient = val.Ip
			return freeClient
		}
	}
	freeClient = ""
	return freeClient
}
