package Minio

import "github.com/minio/minio-go/v7"

var minioClient *minio.Client //Minio连接
var err error                 //错误消息
var BucketName string         //当前统一存储桶
var RawBucketName string      //源文件存储桶

//获取存储桶名
func GetBucket() string {
	return BucketName
}
