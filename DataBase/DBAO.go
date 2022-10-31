package DataBase

import (
	"context"

	"UAutoServer/Logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mong *mongo.Client
var erMainTData []MainTable

type MainTable struct {
	AppKey   string
	UUID     string
	RawFiles []string
}

func InitDB() {
	var err error
	clientOption := options.Client().ApplyURI("mongodb://10.11.145.15:27171")

	//连接到MongoDB
	mong, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		Logs.Error(err)
	}
	//检查连接状态
	err = mong.Ping(context.TODO(), nil)
	if err != nil {
		Logs.Error(err)
	}
	Logs.Print("数据库初始化完毕----")
}
