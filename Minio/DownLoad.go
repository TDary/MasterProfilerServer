package Minio

import (
	"MasterServer/Logs"
	"context"

	"github.com/minio/minio-go/v7"
)

func DownLoadFile(objectName string, filePath string, contentType string) bool {
	downloadMutex.TryLock()
	ctx := context.Background()
	// 检查存储桶是否已经存在。
	exists, err := minioClient.BucketExists(ctx, BucketName)
	if err == nil && exists {
		// Logs.Loggers().Printf("当前存储桶 %s存在----\n", BucketName)
	} else {
		Logs.Loggers().Printf("当前存储桶 %s不存在----\n", BucketName)
		Logs.Loggers().Print(err)
		downloadMutex.Unlock()
		return false
	}
	// 使用FGetObject下载文件。
	err = minioClient.FGetObject(ctx, BucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		Logs.Loggers().Println(err)
		downloadMutex.Unlock()
		return false
	}
	downloadMutex.Unlock()
	return true
}
