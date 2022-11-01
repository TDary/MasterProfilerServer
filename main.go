package main

import (
	"UAutoServer/AnalyzeServer"
	"UAutoServer/DataBase"
	"UAutoServer/HttpServer"
	"UAutoServer/Logs"
)

func main() {
	Logs.Print("Welcome to use UAutoServer")
	DataBase.InitDB()
	AnalyzeServer.InitClient()
	go AnalyzeServer.AnalyzeRequestUrl()
	HttpServer.ListenAndServer("10.11.144.31:8201")
}
