//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"UAutoServer/Logs"
	"UAutoServer/StartSystem"
)

func main() {
	Logs.Loggers().Print("Welcome to use UAutoServerMaster")
	StartSystem.Run()
}
