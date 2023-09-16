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
	UUID      string
	RawFile   string
	AnalyzeIP string
	State     int
}

type InsertSimple struct {
	UUID  string
	Name  string
	Valus []float32
}

type CaseFunRow struct {
	UUID   string
	Name   string
	Frames []FunRowInfo
}

type FunRowInfo struct {
	Frame   int32
	Total   int32
	Self    int32
	Calls   int32
	Gcalloc int32
	Timems  int32
	Selfms  int32
}
