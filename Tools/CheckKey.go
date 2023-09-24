package Tools

import (
	"MasterServer/Logs"
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetKey(k string, m map[int]string) bool {
	for i := 0; i < len(m); i++ {
		if m[i] == k {
			return true
		}
	}
	return false
}

//解压zip文件
func ExtractZip(zipFile string, targetFolder string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		// 获取相对路径
		relPath := strings.TrimPrefix(file.Name, filepath.Dir(file.Name))

		// 拼接目标文件路径
		targetPath := filepath.Join(targetFolder, relPath)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

//发送机器人提醒消息
func SendRobotMsg(url string, msg string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		// handle error
		Logs.Loggers().Print("Failed to send meesage to robot----", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		Logs.Loggers().Print("Failed to send meesage to robot----", err.Error())
	}
	defer resp.Body.Close()
}
