package AnalyzeServer

import (
	"MasterServer/Logs"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func RequestClientState(m_ip string) ReceiveDate {
	//发送开始解析的相关数据信息
	//前提是已经创建好数据表
	request_Url := "http://" + m_ip + "/manastate?testrequest"
	//超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(request_Url)
	if err != nil {
		Logs.Loggers().Print(err.Error())
		return ReceiveDate{}
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
	var rece ReceiveDate
	err = json.Unmarshal(result.Bytes(), &rece)
	if err != nil {
		Logs.Loggers().Print("反序列化失败----", err.Error())
	}
	if rece.Code == 200 {
		//成功获取了
		return rece
	}
	return ReceiveDate{}
}

//发送解析请求
func SendRequestAnalyze(getdata AnalyzeData, m_ip string) {
	//发送开始解析的相关数据信息
	//前提是已经创建好数据表
	request_Url := "http://" + m_ip + "/analyze?uuid=" + getdata.UUID +
		"&rawfile=" + getdata.RawFile + "&rawfilename=" + getdata.RawFileName + "&unityversion=" + getdata.UnityVersion + "&analyzebucket=" + getdata.Bucket
	//超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(request_Url)
	if err != nil {
		Logs.Loggers().Print(err.Error())
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
	if strings.Contains(result.String(), "ok") {
		//客户端成功接收到开始解析的消息，降空闲进程数减1
		Logs.Loggers().Print("解析客户端成功接收到消息，准备开始解析----")
		ReduceRunC(m_ip, 1)
	} else {
		Logs.Loggers().Print("客户端未成功接收到消息----")
		return
	}
}
