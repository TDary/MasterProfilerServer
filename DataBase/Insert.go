package DataBase

import (
	"UAutoServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertMain(data MainTable) {
	//连接数据库表
	col := mong.Database("MyDB").Collection("MainTable")
	//插入数据
	iResult, err := col.InsertOne(context.TODO(), data)
	if err != nil {
		Logs.Print(err)
		erMainTData = append(erMainTData, data)
	}
	//默认生成一个唯一全局ID
	id := iResult.InsertedID.(primitive.ObjectID)
	Logs.Print("插入成功" + id.Hex())
}

func Inserts(datas []MainTable) {
	col := mong.Database("MyDB").Collection("MainTable")
	indata := []interface{}{datas}
	_, err := col.InsertMany(context.TODO(), indata)
	if err != nil {
		Logs.Print(err)
		erMainTData = append(erMainTData, datas...)
	} else {
		Logs.Print("批量插入数据成功")
	}
}
