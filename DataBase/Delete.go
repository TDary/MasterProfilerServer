package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func Del(appkey string, uuid string) {
	col := mong.Database("MyDB").Collection("MainTable")
	many, err := col.DeleteMany(context.TODO(), bson.D{{Key: "uuid", Value: uuid}})
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(many.DeletedCount)
}
