package main

import (
	"UAutoProfiler/AnalyzeServer"
	"fmt"
)

func main() {
	fmt.Print("Welcome to use UAutoServer")
	AnalyzeServer.Run()
	// lastTime := TimeTools.GetLogicTime()
	// for true {
	// 	currentTime := TimeTools.GetLogicTime()
	// 	deltaTime := currentTime.Unix() - lastTime.Unix()
	// 	time.Sleep(1 * time.Second)
	// 	Logs.Print(deltaTime)
	// 	lastTime = currentTime
	// }
}
