package main

import (
	"UAutoServer/Logs"
	"UAutoServer/StartSystem"
)

func main() {
	Logs.Print("Welcome to use UAutoServerMaster")
	StartSystem.Run()
}
