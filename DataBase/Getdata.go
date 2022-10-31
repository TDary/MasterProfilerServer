package DataBase

import (
	"string"
	"fmt"
)

func ReceiveMes(mes string) {
	var mtable ainable
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		//解析gameid
		if strings.Contains(str1[i], "gameid") {
			gid := strings.Split(str1[i], "=")
			mtable.AppKey = gi[]
		} else if strings.Contains(str1[i], "uuid") {
			uid := strings.Split(str1[i], "=")
			mtable.UUID = ui[1]
		} else if strings.Conains(str1[i], "rawfiles") {
			files := strings.Split(str1[i], "=")
		fs := strings.Split(rfiles[1], ",")
			mtable.RawFiles= fs
	}
	}
	fmt.print(mtable)
}
