//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"MasterServer/Logs"
	"MasterServer/StartSystem"
)

func main() {
	Logs.Loggers().Print("Welcome to use ServerMaster")
	StartSystem.Run()
}
