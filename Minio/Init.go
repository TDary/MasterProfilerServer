package Minio

import (
	"MasterServer/Logs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(endpoint string, bucket string, rawbucket string, accessKeyID string, secretAccessKey string) {
	//endpoint //minio服务器url
	BucketName = bucket
	RawBucketName = rawbucket
	useSSL := false
	// 初使化 minio client对象。
	minioClient, err = minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), Secure: useSSL})
	if err != nil {
		Logs.Loggers().Print(err)
		return
	}
	Logs.Loggers().Print("Minio初始化完毕----")
}
