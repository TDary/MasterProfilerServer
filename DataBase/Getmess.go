package DataBase

import (
	"strings"
)

func ReceiveMes(mes string) {
	var mtable MainTable
	str1 := strings.Split(mes, "&")
	for i := 0; i < len(str1); i++ {
		//解析gameid
		if strings.Contains(str1[i], "gameid") {
			gid := strings.Split(str1[i], "=")
			mtable.AppKey = gid[1]
		} else if strings.Contains(str1[i], "uuid") {
			uid := strings.Split(str1[i], "=")
			mtable.UUID = uid[1]
		} else if strings.Contains(str1[i], "rawfiles") {
			files := strings.Split(str1[i], "=")
			fs := strings.Split(files[1], ",")
			mtable.RawFiles = fs
		}
	}
	InsertMain(mtable)
}
