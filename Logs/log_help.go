package Logs

import (
	"log"
	"os"
	"time"
)

var loger *log.Logger

func Init() {
	_, err := os.Stat("./log")
	if err != nil {
		os.Mkdir("./log", 0755)
	}
	file := "./log/" + time.Now().Format("2006-01-02") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[logTool]", log.Ltime|log.Lshortfile)
	// 将文件设置为loger作为输出
}

func Loggers() *log.Logger {
	return loger
}
