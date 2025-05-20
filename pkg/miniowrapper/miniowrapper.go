package miniowrapper

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
)

func GetFilesInPath(ctx context.Context, minioClient *minio.Client, minioBucket string, minioPath string) []string {
	objectsChan := minioClient.ListObjects(ctx, minioBucket, minio.ListObjectsOptions{
		Prefix:    minioPath,
		Recursive: true,
	})

	var files []string
	for object := range objectsChan {
		if object.Err != nil {
			log.Printf("error listing object: %v", object.Err)
			continue
		}

		files = append(files, object.Key)
	}

	return files
}

func DeletePathCascade(ctx context.Context, minioClient *minio.Client, bucket, path string) error {
	objectsChan := minioClient.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: true,
	})

	for object := range objectsChan {
		if object.Err != nil {
			return object.Err
		}

		log.Printf("deleting object: %s", object.Key)
		err := minioClient.RemoveObject(ctx, bucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}
	
	return nil
}