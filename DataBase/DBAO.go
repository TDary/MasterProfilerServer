package DataBase

import (
	"go.mongodb.org/mongo-driver/mongo"
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
	AnalyzeType   string
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
	CsvPath       string
	UnityVersion  string
	AnalyzeBucket string
	AnalyzeIP     string
	StorageIp     string
	State         int
	Priority      string
}
