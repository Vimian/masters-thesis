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

func calculateSizeBytes(dictionaryLength int, windowsAmount int64, windowLengthBytes int64) int64 {
	// math.Floor(math.Log2(len(dictionary)))
	var minimumDictionaryKeyLength int64 = int64(math.Floor(math.Log2(float64(dictionaryLength))))
	if minimumDictionaryKeyLength == 0 {
		minimumDictionaryKeyLength = 1
	}
	//log.Printf("minimum dictionary key length: %d", minimumDictionaryKeyLength)
	
	// math.Ceil(windowsAmount * minimumDictionaryKeyLength / 8)
	var dataSizeBytes int64 = int64(math.Ceil(float64(windowsAmount * minimumDictionaryKeyLength) / 8))
	//log.Printf("data size bytes: %d", dataSizeBytes)

	// math.Ceil(len(dictionary) * windowLengthBytes)
	var dictionarySizeBytes int64 = int64(dictionaryLength) * windowLengthBytes
	//log.Printf("dictionary size bytes: %d", dictionarySizeBytes)

	currentSizeBytes := dataSizeBytes + dictionarySizeBytes
	//log.Printf("current size bytes: %d", currentSizeBytes)
	return currentSizeBytes
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

	var windowLimitInBytes int64 = int64(len(data) / 2)
	var windowLengthBytes int64 = 2
	
	for ; windowLengthBytes <= windowLimitInBytes; windowLengthBytes++ {
		dictionary := make(map[string]bool)
		var dictionaryLimitReached int8 = 0
		
		// math.Ceil(len(data) / windowLengthBytes)
		var windowsAmount int64 = int64(math.Ceil(float64(len(data)) / float64(windowLengthBytes)))
		
		var i int64 = 0
		for ; i < int64(len(data)); i += windowLengthBytes {
			end := i + windowLengthBytes
			if end > int64(len(data)) {
				end = int64(len(data))
			}

			key := string(data[i:end])
			dictionary[key] = true

			currentSizeBytes := calculateSizeBytes(len(dictionary), windowsAmount, windowLengthBytes)

			if currentSizeBytes >= int64(len(data)) {
				log.Println("limit reached")
				
				dictionaryLimitReached = 1
				break
			}
		}

		//log.Printf("dictionary: %v", dictionary)
		
		analysisResult := persistence.AnalysisResult{
			FilePath: filePath,
			FileName: fileName,
			FileSize: fileSize,
			WindowLengthBytes: windowLengthBytes,
			DictionaryLength: int64(len(dictionary)),
			DictionaryLimitReached: dictionaryLimitReached,
			CompressedSizeBytes: calculateSizeBytes(len(dictionary), windowsAmount, windowLengthBytes),
		}
		
		analysisResults = append(analysisResults, analysisResult)
		//log.Printf("analysis result: %v", analysisResult)
	}
	
	return analysisResults, nil
}