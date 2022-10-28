package AnalyzeServer

import (
	"UAutoServer/Logs"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AllProfilerClient struct {
	Ip            string
	WorkerNumbers int
	WorkType      string
}

type ConfigData map[string]interface{}

func Run() {
	InitClient()
	//HttpServer.ListenAndServer("10.11.144.31:8201")
}

func AnalyzeRequestUrl() {

}

func AddAnalyzeClient() {

}

func InitClient() {
	var data, _ = ioutil.ReadFile("./ServerConfig.json")
	var config ConfigData
	var cData AllProfilerClient
	var err = json.Unmarshal(data, &config)
	if err != nil {
		Logs.Error(err)
	}
	res := config["client"]
	for _, test := range res {

	}
	fmt.Print(res)
}
