package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//更新主表状态值
func UpdateMainTable(appkey string, uuid string, rawFiles []string) {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.D{{Key: "$set", Value: bson.D{{Key: "rawFiles", Value: rawFiles}}}}
	//更改元数据
	many, err := col.UpdateMany(context.TODO(), bson.D{{Key: "AppKey", Value: appkey}, {Key: "UUID", Value: uuid}}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	//打印改变了多少
	Logs.Loggers().Print(many.ModifiedCount)
}

func UpdateData() {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.D{{Key: "$set", Value: bson.D{{Key: "AppKey", Value: "tes2"}, {Key: "UUID", Value: "sahsala"}}}}
	//更改元数据
	many, err := col.UpdateMany(context.TODO(), bson.D{{Key: "AppKey", Value: "sasas"}}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	//打印改变了多少
	Logs.Loggers().Print(many.ModifiedCount)
}

//更新子表任务状态
func ModifySub(uuid string, rawfile string, state int) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: state}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}, {Key: "rawfile", Value: rawfile}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}

//更新子表任务状态
func ModifySubOne(objid int, state int) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: state}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: objid}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}

//更新子表成功状态
func UpdateStates(rawfilename string, uuid string, state int, anaip string) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: state}, {Key: "analyzeip", Value: anaip}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}, {Key: "rawfile", Value: rawfilename}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}

//将失败的任务进行重新解析
func FindAndModify(uuid string, rawfile string) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: 0}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}, {Key: "rawfile", Value: rawfile}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}
