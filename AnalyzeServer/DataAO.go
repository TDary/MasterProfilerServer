package AnalyzeServer

import "net"

var config ConfigData
var allAnalyzeClient []ClientState
var stopMsg []EndData
var allconnector []ConnectPool //连接池

type SuccessData struct {
	UUID    string
	IP      string
	RawFile string
}

type ProfilerClient struct {
	Ip            string
	Port          string
	WorkerNumbers int
	WorkType      string
	State         bool
}

type ServerConfig struct {
	Ip   string
	Port string
}
type ConfigData struct {
	Client          []ProfilerClient
	MasterServer    ServerConfig
	MinioServerPath string
	MinioBucket     string
	MinioRawBucket  string
	MergePath       string
	RobotUrl        string
}

type ClientState struct {
	Ip        string
	IpAddress string
	State     string
	Num       int
}

type AnalyzeData struct {
	UUID         string
	AnalyzeType  string
	RawFile      string
	RawFileName  string
	UnityVersion string
	Bucket       string
	Appkey       string
}

type ReceiveDate struct {
	Code  int    `json:"code"`
	State string `json:"state"`
	Num   int    `json:"num"`
}

type EndData struct {
	UUID        string
	LastRawFile string
}

type ConnectPool struct { //以IP为区分，存储连接池对象
	Ip       string
	Marchine string //机器标识
	Conn     net.Conn
}
