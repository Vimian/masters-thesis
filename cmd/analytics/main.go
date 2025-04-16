package main

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vimian/masters-thesis/cmd/analytics/persistence"
	"github.com/vimian/masters-thesis/pkg/miniowrapper"
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

	db, err := persistence.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDatabase)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	defer db.Close()

	if err := persistence.CleanAnalytics(); err != nil {
		log.Fatalf("error cleaning measurements: %v", err)
	}

	// run analytics
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioOriginalPath := os.Getenv("MINIO_ORIGINAL_PATH")

	analytics(minioClient, minioBucket, minioOriginalPath)
}

func analytics(minioClient *minio.Client, minioBucket, minioOriginalPath string) {
	ctx := context.Background()

	filePaths := miniowrapper.GetFilesInPath(ctx, minioClient, minioBucket, minioOriginalPath)

	for _, filePath := range filePaths {
		log.Printf("getting object: %s", filePath)
                objectInfo, err := minioClient.StatObject(ctx, minioBucket, filePath, minio.StatObjectOptions{})
                if err != nil {
                        log.Printf("error stating object %s: %v", filePath, err)
                        continue
                }

                reader, err := minioClient.GetObject(ctx, minioBucket, filePath, minio.GetObjectOptions{})
                if err != nil {
                        log.Printf("error getting object %s: %v", filePath, err)
                        continue
                }
                defer reader.Close()

		analysisResults, err := analyzeFile(reader, filePath, objectInfo)
		if err != nil {
			panic(err)
		}

		for _, analysisResult := range analysisResults {
			if err := persistence.InsertAnalysisResult(analysisResult); err != nil {
				log.Printf("error inserting analysis result: %v", err)
			}
		}
	}
}

func analyzeFile(reader *minio.Object, filePath string, fileInfo minio.ObjectInfo) ([]persistence.AnalysisResult, error) {
	data, err := io.ReadAll(reader)
        if err != nil {
                return nil, err
        }

	var parts []string = strings.Split(filePath, "/")
	var fileName string = parts[len(parts)-1]

	var fileSize int64 = fileInfo.Size

	analysisResults := []persistence.AnalysisResult{}

	var bytesLimit int64 = int64(len(data) / 2)
	var bytes int64 = 2
	
	outer:
	for ; bytes <= bytesLimit; bytes++ {
		dictonary := make(map[string]bool)
		var dictonarySize int64 = 0
		var dictonarySizeLimit int64 = int64(math.Pow(2, float64(bytes - 1))) // TODO: improve to include current size of dictionary
		if dictonarySizeLimit < 0 {
			dictonarySizeLimit = math.MaxInt64
		}
		log.Printf("dictonary size limit: %d", dictonarySizeLimit)

        	buffer := make([]byte, bytes)
		var i int64 = 0
		for ; i < int64(len(data)); i += bytes {
			upperBound := i + bytes
			if upperBound > int64(len(data)) {
				upperBound = int64(len(data)) + 1
			}
			copy(buffer, data[i:upperBound])
			if _, ok := dictonary[string(buffer)]; !ok {
				dictonary[string(buffer)] = true
				dictonarySize++
				log.Printf("dictonary: %v", dictonary)
			}

			if dictonarySize >= dictonarySizeLimit {
				log.Printf("dictonary size limit reached: %d", dictonarySizeLimit)
				log.Printf("dictonary: %v", dictonary)
				
				analysisResult := persistence.AnalysisResult{
					FilePath: filePath,
					FileName: fileName,
					FileSize: fileSize,
					Bytes: bytes,
					BytesNeeded: bytes,
					DictionarySize: dictonarySize,
				}
				analysisResults = append(analysisResults, analysisResult)
				log.Printf("analysis result: %v", analysisResult)
				
				continue outer
			}
		}
		
		var bytesNeeded int64 = 1
		if dictonarySize > 1 {
			bytesNeeded = int64(math.Floor(math.Log2(float64(dictonarySize - 1)))) + 1
		}
		analysisResult := persistence.AnalysisResult{
			FilePath: filePath,
			FileName: fileName,
			FileSize: fileSize,
			Bytes: bytes,
			BytesNeeded: bytesNeeded,
			DictionarySize: dictonarySize,
		}
		analysisResults = append(analysisResults, analysisResult)
		log.Printf("analysis result: %v", analysisResult)
	}
	
	return analysisResults, nil
}
