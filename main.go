package main

import (
	LogC "UAutoProfiler/LogTools"
	GetTimes "UAutoProfiler/TimeTools"
	"fmt"
	"time"
)

func main() {
	fmt.Print("Welcome to use UAutoServer")
	lastTime := GetTimes.GetLogicTime()
	for true {
		currentTime := GetTimes.GetLogicTime()
		deltaTime := currentTime.Unix() - lastTime.Unix()
		time.Sleep(1 * time.Second)
		LogC.Print(deltaTime)
		lastTime = currentTime
	}
}
