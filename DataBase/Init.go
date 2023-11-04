package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(dbAddress string, dbname string, mainTable string, subTable string, funRow string, simpleData string, funPath string) {
	var err error
	clientOption := options.Client().ApplyURI(dbAddress) //"mongodb://192.168.31.40:27171"

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

	databaseName = dbname
	maintable = mainTable
	subtable = subTable
	funrow = funRow
	simpledata = simpleData
	funpath = funPath
	Logs.Loggers().Print("数据库初始化完毕----")
}
