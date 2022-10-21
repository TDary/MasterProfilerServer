package AnalyzeServer

import "UAutoProfiler/HttpServer"

func Run() {
	HttpServer.ListenAndServer("10.11.145.198:8201")
}
