package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

//删除子表数据
func DelSubData(uuid string) {
	col := mong.Database("MyDB").Collection("SubTable")
	_, err := col.DeleteMany(context.TODO(), bson.D{{Key: "uuid", Value: uuid}})
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//删除基础数据
func DelSimpleData(uuid string) {
	col := mong.Database("MyDB").Collection("SimpleData")
	_, err := col.DeleteMany(context.TODO(), bson.D{{Key: "uuid", Value: uuid}})
	if err != nil {
		Logs.Loggers().Print(err)
	}
}

//删除FunRow数据
func DelFunRow(uuid string) {
	col := mong.Database("MyDB").Collection("FunRow")
	_, err := col.DeleteMany(context.TODO(), bson.D{{Key: "uuid", Value: uuid}})
	if err != nil {
		Logs.Loggers().Print(err)
	}
}
