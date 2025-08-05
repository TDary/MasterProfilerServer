package Tools

import (
	"MasterServer/DataBase"
	"MasterServer/Logs"
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 清除数据库无用数据
func ClearDataBase() {
	waitCase := DataBase.FindMainTable(1)
	if len(waitCase) > 0 {
		for _, val := range waitCase {
			DataBase.DelSubData(val.UUID)
		}
	} else {
		Logs.Loggers().Print("无待待删除的子任务数据----")
	}
}

// 解压zip文件
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

// 发送机器人提醒消息
func SendRobotMsg(url string, msg string) {
	var sendArgs strings.Builder
	sendArgs.WriteString(`{"msg_type":"post","content":{"post":{"zh_cn":{"content":[[{"tag":"text","text":"`)
	sendArgs.WriteString(msg)
	sendArgs.WriteString(`"}]]}}}}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(sendArgs.String())))
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

// 使用AES对数据进行解密
func Decrypt(ciphertext, key []byte) ([]byte, error) { //密钥：eb3386a8a8f57a579c93fdfb33ec9471
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 分离Nonce和密文
	nonce := ciphertext[:aes.BlockSize]
	actualCiphertext := ciphertext[aes.BlockSize:]

	// CTR流解密
	stream := cipher.NewCTR(block, nonce)
	plaintext := make([]byte, len(actualCiphertext))
	stream.XORKeyStream(plaintext, actualCiphertext)

	return plaintext, nil
}
