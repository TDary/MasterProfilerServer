package DataBase

import (
	"context"

	"MasterServer/Logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mong *mongo.Client
var erMainTData []MainTable
var erSubTdata []SubTable

type MainTable struct {
	AppKey        string
	UUID          string
	GameName      string
	CaseName      string
	RawFiles      []string
	UnityVersion  string
	AnalyzeBucket string
	StorageIp     string
	Device        string
	TestBeginTime string
	TestEndTime   string
	State         int
	Priority      string
	ScreenState   int
	ScreenFiles   []string
}

type SubTable struct {
	AppKey        string
	UUID          string
	RawFile       string
	UnityVersion  string
	AnalyzeBucket string
	AnalyzeIP     string
	StorageIp     string
	State         int
	Priority      string
}

func InitDB() {
	var err error
	clientOption := options.Client().ApplyURI("mongodb://10.11.145.15:27171")

	//连接到MongoDB
	mong, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	//检查连接状态
	err = mong.Ping(context.TODO(), nil)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	Logs.Loggers().Print("数据库初始化完毕----")
}
