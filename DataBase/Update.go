package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateData() {
	col := mong.Database("MyDB").Collection("MainTable")
	//更改数据
	up := bson.D{{"$set", bson.D{{"AppKey", "tes2"}, {"UUID", "sahsala"}}}}
	//更改元数据
	many, err := col.UpdateMany(context.TODO(), bson.D{{"AppKey", "sasas"}}, up)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	//打印改变了多少
	Logs.Loggers().Print(many.ModifiedCount)
}

//更新子表任务状态
func ModifySub(uuid string, rawfile string, state int) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{"$set", bson.D{{"state", state}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{"uuid", uuid}, {"rawfile", rawfile}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}

//更新子表任务状态
func ModifySubOne(objid int, state int) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{"$set", bson.D{{"state", state}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{"_id", objid}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}

//更新子表成功状态
func UpdateStates(rawfilename string, uuid string, state int, anaip string, csvpath string) {
	col := mong.Database("MyDB").Collection("SubTable")
	update := bson.D{{"$set", bson.D{{"state", state}, {"analyzeiP", anaip}, {"csvpath", csvpath}}}}
	res, err := col.UpdateOne(context.TODO(), bson.D{{"uuid", uuid}, {"rawfile", rawfilename}}, update)
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(res.UpsertedCount)
}
