package DataBase

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var mong *mongo.Client
var erMainTData []MainTable
var erSubTdata []SubTable
var databaseName string
var maintable string
var subtable string
var funrow string
var simpledata string
var funpath string

type MainTable struct {
	AppKey          string
	UUID            string
	GameName        string
	CaseName        string
	RawFiles        []string
	SnapFiles       []string
	UnityVersion    string
	AnalyzeBucket   string
	CollectorIp     string
	AnalyzeType     string
	Device          string
	TestBeginTime   string
	TestEndTime     string
	State           int
	ScreenState     int
	ScreenFiles     []string
	FrameTotalCount int
}

type SubTable struct {
	UUID         string
	RawFile      string
	SnapFile     string
	AnalyzeIP    string
	State        int
	AnalyzeBegin int64
	AnalyzeEnd   int64
}

type InsertSimple struct {
	UUID   string
	Name   string
	Values []float32
}

type CaseFunRowAlone struct {
	UUID         string
	Name         string
	AvgValidTime int32
	Frames       []FunRowInfoAlone
}

type CaseFunRow struct {
	UUID         string
	Name         string
	UnityObject  string
	AvgValidTime int32
	Frames       []FunRowInfo
}

type CaseFunNamePath struct {
	UUID  string
	Stack []string
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

type FunRowInfoAlone struct {
	Frame   int32
	Calls   int32
	Gcalloc int32
	Timems  int32
	Selfms  int32
}
