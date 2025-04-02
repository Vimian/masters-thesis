package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vimian/masters-thesis/cmd/benchmark/persistence"
	"github.com/vimian/masters-thesis/pkg/compression"
	"github.com/vimian/masters-thesis/pkg/compression/exe"
	"github.com/vimian/masters-thesis/pkg/compression/sevenzip"
)

func main() {
	// initialize minio client
	minioServer := os.Getenv("MINIO_SERVER")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioSecure := os.Getenv("MINIO_SECURE") == "true"

	minioClient, err := minio.New(minioServer, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: minioSecure,
	})
	if err != nil {
		log.Fatalf("error creating minio client: %v", err)
	}

        // connect to postgres
        postgresUser := os.Getenv("POSTGRES_USER")
        postgresPassword := os.Getenv("POSTGRES_PASSWORD")
        postgresHost := os.Getenv("POSTGRES_HOST")
        postgresPort := os.Getenv("POSTGRES_PORT")
        postgresDatabase := os.Getenv("POSTGRES_DATABASE")

        db, err := persistence.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDatabase); 
        if err != nil {
                log.Fatalf("error connecting to postgres: %v", err)
        }
        defer db.Close()

        if err := persistence.CleanMeasurements(); err != nil {
                log.Fatalf("error cleaning measurements: %v", err)
        }

	// run benchmark
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioOriginalPath := os.Getenv("MINIO_ORIGINAL_PATH")
	minioCompressedPath := os.Getenv("MINIO_COMPRESSED_PATH")
	minioDecompressedPath := os.Getenv("MINIO_DECOMPRESSED_PATH")

        runs, err := strconv.Atoi(os.Getenv("RUNS"))
        if err != nil {
                log.Fatalf("error parsing runs: %v", err)
        }

        benchmark(minioClient, minioBucket, minioOriginalPath, minioCompressedPath, minioDecompressedPath, runs)
}

func benchmark(minioClient *minio.Client, minioBucket string, minioOriginalPath string, minioCompressedPath string, minioDecompressedPath string, runs int) {
        for i := 0; i < runs; i++ {
                log.Printf("starting run %d", i)
                run(minioClient, minioBucket, minioOriginalPath, minioCompressedPath, minioDecompressedPath, i)
        }

        log.Printf("benchmark complete")
}

func run(minioClient *minio.Client, minioBucket string, minioOriginalPath string, minioCompressedPath string, minioDecompressedPath string, run int) {
        var algorithms = []struct {
                Name      string
                algorithm compression.Algorithm
        }{
                {"PPMd", sevenzip.PPMd},
                {"LZMA", sevenzip.LZMA},
                {"LZMA2", sevenzip.LZMA2},
                {"BZip2", sevenzip.BZip2},
                {"PPMonstr", exe.PPMonstr_exe},
                {"PPMdexe", exe.PPMd_exe},
                {"BitFlipper", compression.BitFlipper{}},
        }

        for i, algorithm := range algorithms {
                ctx := context.Background()

                log.Printf("run %d: algorithm %d: compression", run, i)
                files := getFiles(ctx, minioClient, minioBucket, minioOriginalPath)
                cleanUpPath(ctx, minioClient, minioBucket, minioCompressedPath)
                compressMeasurement := process(ctx, minioClient, minioBucket, minioCompressedPath, algorithm.algorithm.Compress, files)

                log.Printf("run %d: algorithm %d: decompression", run, i)
                files = getFiles(ctx, minioClient, minioBucket, minioCompressedPath)
                cleanUpPath(ctx, minioClient, minioBucket, minioDecompressedPath)
                decompressMeasurement := process(ctx, minioClient, minioBucket, minioDecompressedPath, algorithm.algorithm.Decompress, files)

                if err := saveMeasurement(run, algorithm.Name, compressMeasurement, decompressMeasurement); err != nil {
                        log.Printf("error saving measurement: %v", err)
                }                
        }

        log.Printf("completed run %d successfully", run)
}

func saveMeasurement(run int, algorithmName string, compressMeasurement  []persistence.ProcessMeasurement, decompressMeasurement []persistence.ProcessMeasurement) error {
        if len(compressMeasurement) != len(decompressMeasurement) {
                return fmt.Errorf("different number of files: %d != %d", len(compressMeasurement), len(decompressMeasurement))
        }

        for i := 0; i < len(compressMeasurement); i++ {
                calculateDuration(&compressMeasurement[i])
                calculateDuration(&decompressMeasurement[i])

                measurement := persistence.Measurement{
                        Run:        run,
                        Algorithm:  algorithmName,

                        OriginalPath:    compressMeasurement[i].OriginalPath,
                        OriginalSize:    compressMeasurement[i].OriginalSize,
                        OriginalHash:    compressMeasurement[i].OriginalHash,

                        Compress: compressMeasurement[i],
                        Decompress: decompressMeasurement[i],

                        CompressionRatio: float64(compressMeasurement[i].OriginalSize) / float64(compressMeasurement[i].Size),
                }
                if err := persistence.InsertMeasurement(measurement); err != nil {
                        return err
                }
        }

        return nil
}

func calculateDuration(measurement *persistence.ProcessMeasurement) {
        measurement.DurationGetObjectInfo = measurement.EndGetObjectInfo - measurement.StartGetObjectInfo
        measurement.DurationGetObject = measurement.EndGetObject - measurement.StartGetObject
        measurement.DurationAlgorithm = measurement.EndAlgorithm - measurement.StartAlgorithm
        measurement.DurationUpload = measurement.EndUpload - measurement.StartUpload
}

func getFiles(ctx context.Context, minioClient *minio.Client, minioBucket string, minioPath string) []string {
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

func process(ctx context.Context, minioClient *minio.Client, minioBucket, minioDestinationPath string, algorihtmFunc compression.AlgorithmFunc, files []string) []persistence.ProcessMeasurement {
        processMeasurements := make([]persistence.ProcessMeasurement, len(files))

        for i, sourcePath := range files {
                processMeasurements[i].OriginalPath = sourcePath

                log.Printf("getting object: %s", sourcePath)
                processMeasurements[i].StartGetObjectInfo = time.Now().UnixNano()
                objectInfo, err := minioClient.StatObject(ctx, minioBucket, sourcePath, minio.StatObjectOptions{})
                if err != nil {
                        log.Printf("error stating object %s: %v", sourcePath, err)
                        continue
                }
                processMeasurements[i].EndGetObjectInfo = time.Now().UnixNano()
                processMeasurements[i].OriginalSize = objectInfo.Size
                processMeasurements[i].OriginalHash = objectInfo.ETag

                processMeasurements[i].StartGetObject = time.Now().UnixNano()
                reader, err := minioClient.GetObject(ctx, minioBucket, sourcePath, minio.GetObjectOptions{})
                if err != nil {
                        log.Printf("error getting object %s: %v", sourcePath, err)
                        continue
                }
                defer reader.Close()
                processMeasurements[i].EndGetObject = time.Now().UnixNano()

		log.Printf("running algorithm on object: %s", sourcePath)
                processMeasurements[i].StartAlgorithm = time.Now().UnixNano()
		compressedReader, fileName, err := algorihtmFunc(reader, objectInfo)
		if err != nil {
			log.Printf("error running algorithm on %s: %v", sourcePath, err)
			continue
		}
                processMeasurements[i].EndAlgorithm = time.Now().UnixNano()

                newSize, err := getFileSize(compressedReader)
                if err != nil {
                        log.Printf("error getting file size: %v", err)
                        continue
                }
                processMeasurements[i].Size = newSize

                destinationPath := minioDestinationPath + "/" + fileName
                processMeasurements[i].Path = destinationPath

		log.Printf("uploading object: %s", destinationPath)
                processMeasurements[i].StartUpload = time.Now().UnixNano()
                _, err = minioClient.PutObject(ctx, minioBucket, destinationPath, compressedReader, newSize, minio.PutObjectOptions{
                        ContentType:    objectInfo.ContentType,
                        UserMetadata:   objectInfo.UserMetadata,
                        UserTags:       objectInfo.UserTags,
                })
                if err != nil {
                        log.Printf("error uploading %s: %v", destinationPath, err)
                        continue
                }
                processMeasurements[i].EndUpload = time.Now().UnixNano()

                log.Printf("successfully converted: %s -> %s", sourcePath, destinationPath)
        }

        return processMeasurements
}

func getFileSize(reader io.Reader) (int64, error) {
        if seekable, ok := reader.(io.Seeker); ok {
                _, err := seekable.Seek(0, io.SeekStart)
                if err != nil {
                        log.Printf("error seeking to beginning of stream %v", err)
                        return 0, err
                }

                if size, ok := reader.(interface{ Size() int64 }); ok {
                        return size.Size(), nil
                }
        }

        compressedSize, err := io.Copy(io.Discard, reader)
        if err != nil {
                log.Printf("error determining compressed size %v", err)
                return 0, err
        }

        _, err = reader.(io.Seeker).Seek(0, io.SeekStart)
        if err != nil {
                log.Printf("error seeking to beginning of stream %v", err)
                return 0, err
        }

        return compressedSize, nil
}