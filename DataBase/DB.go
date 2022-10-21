package Database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string
	Value int
}

////use admin
//db.createUser({
//  user: 'root',          // 用户名
//  pwd: 'cdr123',      // 密码
//  roles:[{
//    role: 'root',  // 读写权限角色
//    db: 'admin'     // 数据库名
//  }]
//})
// use MyDB
// db.createUser({
//   user: 'Dary',          // 用户名
//   pwd: 'cdr123',      // 密码
//   roles:[{
//     role: 'root',  // 读写权限角色
//     db: 'MyDB'     // 数据库名
//   }]
// })

func TestConn() {
	var result User
	//设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://Dary:cdr123@10.11.145.15:27171/?authSource=MyDB")

	//连接到Mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		loger.Fatal(err)
	}

	col := client.Database("MyDB").Collection("test")
	err = col.FindOne(context.TODO(), bson.D{{}}).Decode(&result)
	if err != nil {
		loger.Fatal(err)
	}
	loger.Print(result)

	//检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		loger.Fatal(err)
	}
	loger.Println("Connectd to MongoDB!")

}
