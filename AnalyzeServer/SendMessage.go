package AnalyzeServer

import (
	"MasterServer/Logs"
	"encoding/json"
	"net/http"
	"time"
)

func RequestClientState(m_ip string) ReceiveDate {
	//发送开始解析的相关数据信息
	//前提是已经创建好数据表
	request_Url := "manastate?testrequest"
	//超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(request_Url)
	if err != nil {
		Logs.Loggers().Print(err.Error())
		return ReceiveDate{}
	}
	defer resp.Body.Close()
	// var buffer [512]byte
	// result := bytes.NewBuffer(nil)
	var result []byte
	resp.Body.Read(result)
	var rece ReceiveDate
	err = json.Unmarshal(result, &rece)
	if err != nil {
		Logs.Loggers().Print("反序列化失败----", err.Error(), string(result))
	}
	if rece.Code == 200 {
		//成功获取了
		return rece
	}
	return ReceiveDate{}
}

//通过socket发送解析请求
func SendRequestAnalyze(getdata AnalyzeData, ip string) {
	requestMsg := "analyze?uuid=" + getdata.UUID +
		"&rawfile=" + getdata.RawFile + "&rawfilename=" + getdata.RawFileName + "&unityversion=" + getdata.UnityVersion + "&analyzebucket=" + getdata.Bucket
	GetConn(ip).Write([]byte(requestMsg))
}
