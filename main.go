package main

import (
	"UAutoServer/Logs"
	"UAutoServer/StartSystem"
)

func main() {
	Logs.Loggers().Print("Welcome to use UAutoServerMaster")
	StartSystem.Run()
}
