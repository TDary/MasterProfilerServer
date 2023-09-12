package Minio

import (
	"MasterServer/Logs"
	"context"

	"github.com/minio/minio-go/v7"
)

//查询存储桶中的对象
func SearchObjectOfBucket(uuid string) []string {
	var allrawfile []string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := minioClient.ListObjects(ctx, RawBucketName, minio.ListObjectsOptions{
		WithVersions: true,
		WithMetadata: false,
		Prefix:       uuid,
		Recursive:    true,
	})
	for object := range objectCh {
		if object.Err != nil {
			Logs.Loggers().Println(object.Err)
			return nil
		}
		ishasSeem := false
		for _, val := range allrawfile {
			if val == object.Key {
				//有同项，不新增
				ishasSeem = true
			}
			if !ishasSeem { //没同项，新增
				allrawfile = append(allrawfile, object.Key)
			}
		}
	}
	return allrawfile
}
