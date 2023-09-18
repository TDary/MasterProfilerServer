package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//更新主表
func ModifyMain(uuid string, state int, framecount int) {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.M{"$set": bson.M{"state": state, "frametotalcount": framecount}}
	//更改元数据
	_, err := col.UpdateMany(context.TODO(), bson.M{"uuid": uuid}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//更新主表状态值
func ModifyMainState(uuid string, state int) {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.M{"$set": bson.M{"state": state}}
	//更改元数据
	_, err := col.UpdateMany(context.TODO(), bson.M{"uuid": uuid}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//更新主表
func UpdateMainTable(appkey string, uuid string, rawFiles []string) {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.M{"$set": bson.M{"rawfiles": rawFiles}}
	//更改元数据
	_, err := col.UpdateMany(context.TODO(), bson.M{"appkey": appkey, "uuid": uuid}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//更新子表成功状态
func UpdateStates(rawfilename string, uuid string, state int, anaip string) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.M{"$set": bson.M{"state": state, "analyzeip": anaip}}
	_, err := col.UpdateOne(context.TODO(), bson.M{"uuid": uuid, "rawfile": rawfilename}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//将失败的任务进行重新解析
func FindAndModify(uuid string, rawfile string) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: 0}}}}
	_, err := col.UpdateOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}, {Key: "rawfile", Value: rawfile}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}
