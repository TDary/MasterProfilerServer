package AnalyzeServer

var config ConfigData
var isAnalyzeStop bool                    //完成解析处理控制信号
var allclients map[string]*ProfilerClient //解析客户端及服务端配置
var allAnalyzeClient []ClientState
var stopMsg []EndData

type SuccessData struct {
	UUID    string
	IP      string
	RawFile string
	CsvPath string
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
	MergeServer     ServerConfig
	MasterServer    ServerConfig
	MinioServerPath string
	MinioBucket     string
	MinioRawBucket  string
}

type ClientState struct {
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
	Ip          string
	UUID        string
	LastRawFile string
}
