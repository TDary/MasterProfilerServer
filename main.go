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

//服务器被强行关机的情况要做处理!!!!!!!断线重连的情况等
