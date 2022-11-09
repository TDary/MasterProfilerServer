package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertMain(data MainTable) {
	//连接数据库表
	col := mong.Database("MyDB").Collection("MainTable")
	//插入数据
	iResult, err := col.InsertOne(context.TODO(), data)
	if err != nil {
		Logs.Loggers().Print(err)
		erMainTData = append(erMainTData, data)
	}
	//默认生成一个唯一全局ID
	id := iResult.InsertedID.(primitive.ObjectID)
	Logs.Loggers().Print("插入成功" + id.Hex())
}

func InsertSub(data SubTable) {
	col := mong.Database("MyDB").Collection("SubTable")
	Result, err := col.InsertOne(context.TODO(), data)
	if err != nil {
		Logs.Loggers().Print(err)
		erSubTdata = append(erSubTdata, data)
	}
	id := Result.InsertedID.(primitive.ObjectID)
	Logs.Loggers().Print("插入成功" + id.Hex())
}

func InsertsMain(datas []MainTable) {
	col := mong.Database("MyDB").Collection("MainTable")
	indata := []interface{}{datas}
	_, err := col.InsertMany(context.TODO(), indata)
	if err != nil {
		Logs.Loggers().Print(err)
		erMainTData = append(erMainTData, datas...)
	} else {
		Logs.Loggers().Print("批量插入数据成功----")
	}
}

func InsertsSub(datas []SubTable) {
	col := mong.Database("MyDB").Collection("SubTable")
	indata := []interface{}{datas}
	_, err := col.InsertMany(context.TODO(), indata)
	if err != nil {
		Logs.Loggers().Print(err)
		erSubTdata = append(erSubTdata, datas...)
	} else {
		Logs.Loggers().Print("批量插入子表数据成功----")
	}
}
