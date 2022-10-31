package main

import (
	"UAutoServer/AnalyzeServer"
	"UAutoServer/Logs"
)

func main() {
	Logs.Print("Welcome to use UAutoServer")
	AnalyzeServer.Run()
}
