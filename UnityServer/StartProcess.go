package UServer

import (
	LogC "UAutoProfiler/LogTools"
	"os/exec"
)

//启动解析进程
func StartAnalyze() {
	Unity_Name := "E:/Unity/2020.3.0f1c1/Editor/Data/2020.3.37f1c1/Editor/Unity.exe " //需要启动的程序名,Unity.exe的具体目录
	argu := "-quit -batchmode -nographics "
	argu = argu + "-projectPath E:/U3DProfiler/U3D_ProfilerSDK "
	argu = argu + "-executeMethod Entrance.EntranceParseBegin "
	argu = argu + "-logFile E:/Result/test.log "
	argu = argu + "-rawPath E:/Result/test.raw "
	argu = argu + "-csvPath E:/Result/test.csv "
	argu = argu + "-funjsonPath E:/Result/testfunjson.json "
	argu = argu + "-funrowjsonPath E:/Result/testfunrow.json "
	argu = argu + "-funrenderrowjsonPath E:/Result/testrenderrow.json "
	argu = argu + "-funhashPath E:/Result/testfunhash.json "
	argu = argu + "-Index 0 "
	argu = argu + "-shieldSwitch false "
	argu = Unity_Name + argu
	cmd := exec.Command("cmd.exe", "/c", "start "+argu)
	er := cmd.Run()
	if er != nil { // 运行命令
		LogC.Print(er)
	}
}
