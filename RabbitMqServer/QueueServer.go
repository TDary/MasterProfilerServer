package RabbitMqServer

import (
	"MasterServer/Logs"
	"MasterServer/RabbitMqServer/queue"
	"context"
	"sync"
)

var mutex sync.Mutex

// 推送数据
func PutData(tmp string, msg string) {
	mutex.TryLock()
	q, err := queue.NewFifoDiskQueue("./" + tmp)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	defer q.Close()
	_ = q.Put(context.Background(), []byte(msg))
	mutex.Unlock()
}

//获取数据,每读取一次就会将当前的数据在队列中删除
func GetData(tmp string) string {
	mutex.TryLock()
	q, err := queue.NewFifoDiskQueue("./" + tmp)
	if err != nil {
		Logs.Loggers().Fatal(err)
	}
	defer q.Close()

	result, err := q.Get(context.TODO())
	if err != nil && err == queue.ErrQueueEmpty {
		//Logs.Loggers().Print(err)
		mutex.Unlock()
		return ""
	}
	mutex.Unlock()
	return string(result)
}
