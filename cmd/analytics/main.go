package main

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vimian/masters-thesis/cmd/analytics/algorithms"
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

		analysisResultsStaticWindow, err := algorithms.StaticWindow{}.AnalyzeFile(reader, filePath, objectInfo)
		if err != nil {
			panic(err)
		}

		for _, analysisResult := range analysisResultsStaticWindow {
			if err := persistence.InsertAnalysisResult(analysisResult); err != nil {
				log.Printf("error inserting analysis result: %v", err)
			}
		}

		/*analysisResultsDynamicWindow, err := algorithms.DynamicWindow{}.AnalyzeFile(reader, filePath, objectInfo)
		if err != nil {
			panic(err)
		}
		return

		for _, analysisResult := range analysisResultsDynamicWindow {
			if err := persistence.InsertAnalysisResult(analysisResult); err != nil {
				log.Printf("error inserting analysis result: %v", err)
			}
		}

		return*/
	}
}
