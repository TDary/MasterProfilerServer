package Logs

import (
	"log"
	"os"
	"time"
)

var loger *log.Logger

func init() {
	file := "./log/" + time.Now().Format("2006-01-02") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[logTool]", log.LstdFlags|log.Lshortfile|log.LUTC)
	// 将文件设置为loger作为输出
	return
}

func Print(message any) {
	loger.Print(message)
}

//使用于强制结束进程,到此会直接关闭服务进程
func Error(message any) {
	loger.Fatal(message)
}
