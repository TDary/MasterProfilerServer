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

//更新解析失败的任务状态
func UpdatSubTableFailedStates(rawfilename string, uuid string, state int) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.M{"$set": bson.M{"state": state, "analyzebegin": 0}}
	_, err := col.UpdateOne(context.TODO(), bson.M{"uuid": uuid, "rawfile": rawfilename}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//更新子表成功状态
func UpdateSuccessStates(rawfilename string, uuid string, state int, anaip string, unixTime int64) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.M{"$set": bson.M{"state": state, "analyzeip": anaip, "analyzeend": unixTime}}
	_, err := col.UpdateOne(context.TODO(), bson.M{"uuid": uuid, "rawfile": rawfilename}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//将失败的任务进行重新解析
func FindAndModify(uuid string, rawfile string, state int, unixTime int64) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.M{"$set": bson.M{"state": state, "analyzebegin": unixTime}}
	_, err := col.UpdateOne(context.TODO(), bson.M{"uuid": uuid, "rawfile": rawfile}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
}
