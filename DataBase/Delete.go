package DataBase

import (
	"MasterServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func Del() {
	col := mong.Database("MyDB").Collection("MainTable")
	many, err := col.DeleteMany(context.TODO(), bson.D{{"AppKey", "test"}})
	if err != nil {
		Logs.Loggers().Print(err)
	}
	Logs.Loggers().Print(many.DeletedCount)
}
