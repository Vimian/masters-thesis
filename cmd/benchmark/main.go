package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vimian/masters-thesis/pkg/compression"
)

func main() {
	minioServer := os.Getenv("MINIO_SERVER")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioSecure := os.Getenv("MINIO_SECURE") == "true"

	minioBucket := os.Getenv("MINIO_BUCKET")
	minioOriginalPath := os.Getenv("MINIO_ORIGINAL_PATH")
	minioCompressedPath := os.Getenv("MINIO_COMPRESSED_PATH")
	minioDecompressedPath := os.Getenv("MINIO_DECOMPRESSED_PATH")

	// initialize minio client
	minioClient, err := minio.New(minioServer, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: minioSecure,
	})
	if err != nil {
		log.Fatalf("error creating minio client: %v", err)
	}

	// run compression benchmark
	compressFiles(minioClient, minioBucket, minioOriginalPath, minioCompressedPath)

	// run decompress benchmark
	decompressFiles(minioClient, minioBucket, minioCompressedPath, minioDecompressedPath)
}

func cleanUpPath(ctx context.Context, minioClient *minio.Client, bucket, path string) error {
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

func process(minioClient *minio.Client, minioBucket, minioSourcePath, minioDestinationPath string, algorihtmFunc compression.AlgorithmFunc) {
	ctx := context.Background()

	err := cleanUpPath(ctx, minioClient, minioBucket, minioDestinationPath)
        if err != nil {
                log.Fatalf("error cleaning up path %s: %v", minioDestinationPath, err)
        }

        objectsChan := minioClient.ListObjects(ctx, minioBucket, minio.ListObjectsOptions{
		Prefix:    minioSourcePath,
                Recursive: true,
        })

        for object := range objectsChan {
                if object.Err != nil {
                        log.Printf("error listing object: %v", object.Err)
                        continue
                }

                sourcePath := object.Key

                log.Printf("getting object: %s", sourcePath)

                objectInfo, err := minioClient.StatObject(ctx, minioBucket, sourcePath, minio.StatObjectOptions{})
                if err != nil {
                        log.Printf("error stating object %s: %v", sourcePath, err)
                        continue
                }

                reader, err := minioClient.GetObject(ctx, minioBucket, sourcePath, minio.GetObjectOptions{})
                if err != nil {
                        log.Printf("error getting object %s: %v", sourcePath, err)
                        continue
                }
                defer reader.Close()

		log.Printf("running algorithm on object: %s", sourcePath)

		compressedReader, fileName, err := algorihtmFunc(reader, objectInfo)
		if err != nil {
			log.Printf("error running algorithm on %s: %v", sourcePath, err)
			continue
		}

                // Calculate the size of the compressed object
                var compressedSize int64
                if seekable, ok := compressedReader.(io.Seeker); ok {
                        _, err := seekable.Seek(0, io.SeekStart) //rewind
                        if err != nil {
                                log.Printf("Error seeking to beginning of compressed stream %v", err)
                                continue //skip to the next file
                        }

                        if size, ok := compressedReader.(interface{ Size() int64 }); ok {
                                compressedSize = size.Size()
                        } else { // Handle the case where Size() method is not available. Read until io.EOF
                                var err error
                                compressedSize, err = io.Copy(io.Discard, compressedReader)
                                if err != nil {
                                        log.Printf("Error determining compressed size %v", err)
                                        continue //skip to the next file
                                }
                        }
                } else {
                        var err error
                        compressedSize, err = io.Copy(io.Discard, compressedReader)
                        if err != nil {
                                log.Printf("Error determining compressed size %v", err)
                                continue //skip to the next file
                        }

                }

                _, err = compressedReader.(io.Seeker).Seek(0, io.SeekStart)
                if err != nil {
                        log.Printf("Error seeking to beginning of compressed stream %v", err)
                        continue
                }

                destinationPath := minioDestinationPath + "/" + fileName

		log.Printf("uploading object: %s", destinationPath)
                
                _, err = minioClient.PutObject(ctx, minioBucket, destinationPath, compressedReader, compressedSize, minio.PutObjectOptions{
                        ContentType:    objectInfo.ContentType,
                        UserMetadata:   objectInfo.UserMetadata,
                        UserTags:       objectInfo.UserTags,
                })
                if err != nil {
                        log.Printf("error uploading %s: %v", destinationPath, err)
                        continue
                }

                log.Printf("success: %s -> %s", sourcePath, destinationPath)
        }
}

func compressFiles(minioClient *minio.Client, minioBucket, minioOriginalPath, minioCompressedPath string) {
	log.Printf("running compression benchmark")

	process(minioClient, minioBucket, minioOriginalPath, minioCompressedPath, compression.PPMd{}.Compress)
}

func decompressFiles(minioClient *minio.Client, minioBucket, minioCompressedPath, minioDecompressedPath string) {
	log.Printf("running decompression benchmark")

	process(minioClient, minioBucket, minioCompressedPath, minioDecompressedPath, compression.PPMd{}.Decompress)
}