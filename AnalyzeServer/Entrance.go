package AnalyzeServer

import "UAutoServer/HttpServer"

func Run() {
	HttpServer.ListenAndServer("10.11.144.31:8201")
}
