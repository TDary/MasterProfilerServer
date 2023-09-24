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

//插入子表数据
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

//插入基础数据
func InsertSimpleData(datas []InsertSimple) {
	// 将 []InsertSimple 转换为 []interface{}
	// 创建新的 []interface{} 切片，并从 []InsertSimple 复制值
	interfaceSlice := make([]interface{}, len(datas))
	for i, v := range datas {
		interfaceSlice[i] = v
	}
	col := mong.Database("MyDB").Collection("SimpleData")
	_, err := col.InsertMany(context.TODO(), interfaceSlice)
	if err != nil {
		Logs.Loggers().Print(err)
	} else {
		Logs.Loggers().Print("SimpleData插入数据成功----")
	}
}

//插入FunRow数据
func InsertCaseFunRow(datas []CaseFunRow) {
	// 将 []CaseFunRow 转换为 []interface{}
	// 创建新的 []interface{} 切片，并从 []CaseFunRow 复制值
	interfaceSlice := make([]interface{}, len(datas))
	for i, v := range datas {
		interfaceSlice[i] = v
	}
	col := mong.Database("MyDB").Collection("FunRow")
	_, err := col.InsertMany(context.Background(), interfaceSlice)
	if err != nil {
		Logs.Loggers().Print(err.Error())
	} else {
		Logs.Loggers().Print("FunRow插入数据成功----")
	}
}

func InsertFunNamePath(datas CaseFunNamePath) {
	col := mong.Database("MyDB").Collection("FunNamePath")
	_, err := col.InsertOne(context.Background(), datas)
	if err != nil {
		Logs.Loggers().Print(err.Error())
	} else {
		Logs.Loggers().Print("FunNamePath插入数据成功----")
	}
}
