package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func FindMainTable(state int) []MainTable {
	filter := bson.D{{Key: "state", Value: state}}
	col := mong.Database("MyDB").Collection("MainTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败----", err.Error())
		return nil
	}
	var MainT []MainTable
	err = res.All(context.TODO(), &MainT)
	if err != nil {
		Logs.Loggers().Print("查询失败----", err.Error())
	}
	return MainT
}

func FindSubTableData(uuid string) []SubTable {
	filter := bson.D{{Key: "uuid", Value: uuid}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败:", err)
		return nil
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	if err != nil {
		Logs.Loggers().Print("查询失败----", err.Error())
	}
	return resSubT
}

//查找正常解析案例
func FindSTbyState(state int) []SubTable {
	filter := bson.D{{Key: "state", Value: state}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败:", err)
		return nil
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	if err != nil {
		Logs.Loggers().Print("查询失败----", err.Error())
	}
	return resSubT
}

//查找有高优先级的案例
func FindSTHigh(state int) []SubTable {
	filter := bson.D{{Key: "state", Value: state}, {Key: "priority", Value: "high"}}
	col := mong.Database("MyDB").Collection("SubTable")
	res, err := col.Find(context.TODO(), filter)
	if err != nil {
		Logs.Loggers().Println("查询失败:", err)
	}
	var resSubT []SubTable
	err = res.All(context.TODO(), &resSubT)
	if err != nil {
		Logs.Loggers().Print("查询失败----", err.Error())
	}
	return resSubT
}
