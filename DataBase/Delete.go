package DataBase

import (
	"UAutoServer/Logs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func Del() {
	col := mong.Database("MyDB").Collection("MainTable")
	many, err := col.DeleteMany(context.TODO(), bson.D{{"AppKey", "test"}})
	if err != nil {
		Logs.Print(err)
	}
	Logs.Print(many.DeletedCount)
}
