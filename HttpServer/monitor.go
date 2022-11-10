package HttpServer

import (
	"MasterServer/AnalyzeServer"
	"encoding/json"
	"net/http"
	"strings"
)

func ListenAndServer(address string) {
	http.HandleFunc("/RequestProfiler", RequestProfiler)
	http.HandleFunc("/SuccessProfiler", SuccessProfiler)
	http.HandleFunc("/RquestClient", RequestClient)
	http.HandleFunc("/redirect", Redirect)
	//Http监听函数
	http.ListenAndServe(address, nil)
}

//Http请求处理模块
func DealReceivedMessage(msg string) int {
	if strings.Contains(msg, "RequestProfiler") {
		//Http://serverip:port/RequestProfiler?gameid=test&uuid=test&rawfiles=1.raw,2.raw,3.raw
		beginMsg := strings.Split(msg, "?")[1]
		go StorageParseMes(beginMsg)
		return 200
	} else if strings.Contains(msg, "SuccessProfiler") {
		suce := strings.Split(msg, "?")[1]
		go StorageSucessParseMes(suce)
		return 200
	} else if strings.Contains(msg, "RquestClient") {
		req := strings.Split(msg, "?")[1]
		go AnalyzeServer.AddAnalyzeClient(req)
		return 200
	} else {
		return 400
		//TODO:扩展处理模块
	}
}

//请求解析响应模块
func RequestProfiler(w http.ResponseWriter, r *http.Request) {
	var resData string
	RequestUrlData := r.URL.String()
	resMes := DealReceivedMessage(RequestUrlData)
	if resMes == 200 {
		resData = "success"
	} else {
		resData = "Request Fail"
	}
	w.Header().Set("Content-Type", "application/json") //设置响应内容
	res := Res{
		Code: resMes,
		Data: resData,
	}
	jsonByte, _ := json.Marshal(res) //转json
	w.Write(jsonByte)
}

//解析成功消息回调处理
func SuccessProfiler(w http.ResponseWriter, r *http.Request) {
	var resData string
	RequestUrlData := r.URL.String()
	resMes := DealReceivedMessage(RequestUrlData)
	if resMes == 200 {
		resData = "success"
	} else {
		resData = "Request Fail"
	}
	w.Header().Set("Content-Type", "application/json") //设置响应内容
	res := Res{
		Code: resMes,
		Data: resData,
	}
	jsonByte, _ := json.Marshal(res) //转json
	w.Write(jsonByte)
}

//请求并入组网
func RequestClient(w http.ResponseWriter, r *http.Request) {
	var resData string
	RequestUrlData := r.URL.String()
	resMes := DealReceivedMessage(RequestUrlData)
	if resMes == 200 {
		resData = "success"
	} else {
		resData = "Request Fail"
	}
	w.Header().Set("Content-Type", "application/json") //设置响应内容
	res := Res{
		Code: resMes,
		Data: resData,
	}
	jsonByte, _ := json.Marshal(res) //转json
	w.Write(jsonByte)
}

//重定向功能模块(测试中，可添加其他功能)
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Localtion", "https://www.baidu.com")
	w.WriteHeader(302) //设置响应状态码
}
