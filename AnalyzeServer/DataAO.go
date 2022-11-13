package AnalyzeServer

var config ConfigData
var isStop bool                           //请求解析处理控制信号
var isAnalyzeStop bool                    //完成解析处理控制信号
var allclients map[string]*ProfilerClient //解析客户端及服务端配置

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

type MergeServerConfig struct {
	Ip   string
	Port string
}
type ConfigData struct {
	Client      []ProfilerClient
	MergeServer MergeServerConfig
}
