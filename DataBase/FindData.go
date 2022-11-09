package DataBase

import (
	"UAutoServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func FindSubTableData(uuid string) []SubTable {
	filter := bson.D{{"uuid", uuid}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Print("查询失败UUID：" + uuid)
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	return resSubT
}

//查找正常解析案例
func FindSTbyState(state int) []SubTable {
	filter := bson.D{{"state", state}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败State：", state)
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	return resSubT
}

//查找有高优先级的案例
func FindSTHigh(state int) []SubTable {
	filter := bson.D{{"state", state}, {"priority", "high"}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败State：", state)
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	return resSubT
}
