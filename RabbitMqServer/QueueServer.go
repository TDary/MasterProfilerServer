package RbQServer

import (
	"UAutoServer/Logs"
	"UAutoServer/queue"
	"context"
)

// 推送数据
func PutData(tmp string, msg string) {
	q, err := queue.NewFifoDiskQueue("./" + tmp)
	if err != nil {
		Logs.Error(err)
	}
	defer q.Close()
	_ = q.Put(context.Background(), []byte(msg))
}

//获取数据,每读取一次就会将当前的数据在队列中删除
func GetData(tmp string) string {
	q, err := queue.NewFifoDiskQueue("./" + tmp)
	if err != nil {
		Logs.Error(err)
	}
	defer q.Close()

	result, err := q.Get(context.TODO())
	if err != nil && err == queue.ErrQueueEmpty {
		Logs.Print(err)
		return ""
	}

	return string(result)
}

//删除数据
func DeleteData(tmp string) {
	//TODO:删除数据库中的数据，停止缓存写入队列
}
