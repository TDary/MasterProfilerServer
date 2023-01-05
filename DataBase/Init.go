package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() {
	var err error
	clientOption := options.Client().ApplyURI("mongodb://10.11.144.31:27171")

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
