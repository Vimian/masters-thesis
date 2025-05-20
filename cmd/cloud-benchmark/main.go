package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/vimian/masters-thesis/cmd/cloud-benchmark/persistence"
)

func main() {
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

	if err := persistence.CleanCloudResults(); err != nil {
		log.Fatalf("error cleaning measurements: %v", err)
	}

	filesnames := getFileNames("./files")
	if filesnames == nil {
		log.Fatalf("error getting file names")
	}

	numberOfRuns, err := strconv.Atoi(os.Getenv("RUNS"))
	if err != nil {
		log.Fatalf("error parsing runs: %v", err)
	}
	test(numberOfRuns, filesnames)
	fmt.Println("All tests completed successfully.")
}

func getFileNames(directory string) []string {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	filenames := []string{}
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func test(numberOfRuns int, fileNames []string) {
	//premium do not work with the current solution, don't know if it can?
	//maPremiumKey := os.Getenv("MA_PREMIUM_TIER")
	maHotKey := os.Getenv("MA_HOT_TIER")
	maCoolKey := os.Getenv("MA_COOL_TIER")
	maColdKey := os.Getenv("MA_COLD_TIER")

	for i := 0; i < numberOfRuns; i++ {
		for _, fileName := range fileNames {
			run(i, "hot", maHotKey, fileName)
		}
		for _, fileName := range fileNames {
			run(i, "cool", maCoolKey, fileName)
		}
		for _, fileName := range fileNames {
			run(i, "cold", maColdKey, fileName)
		}
	}
}

func run(run int, tierName string, connectionString string, fileName string) {
	containerName := "uploadcontainer"
	localFilePath := "./files/" + fileName
	downloadedFilePath := "./downloads/" + fileName

	// Create a BlobServiceClient
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		fmt.Println("Error creating service client:", err)
		return
	}

	// Upload the file
	fmt.Println("Starting file upload...")

	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileSize := fileInfo.Size()

	// Hash the original file
	originalHash := sha256.New()
	if _, err := io.Copy(originalHash, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}
	fmt.Println("original hash:", hex.EncodeToString(originalHash.Sum(nil)))

	uploadStartTime := time.Now().UnixNano()

	//upload the file to the blob storage
	_, err = client.UploadFile(context.Background(), containerName, fileName, file, nil)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	uploadEndTime := time.Now().UnixNano()
	uploadDuration := uploadEndTime - uploadStartTime
	fmt.Printf("File uploaded successfully to '%s' in container '%s'.\n", fileName, containerName)

	// Download the file
	fmt.Println("\nStarting file download...")

	downloadFile, err := os.Create(downloadedFilePath)
	if err != nil {
		fmt.Println("Error creating local file:", err)
		return
	}
	defer downloadFile.Close()

	downloadStartTime := time.Now().UnixNano()

	_, err = client.DownloadFile(context.Background(), containerName, fileName, downloadFile, nil)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	downloadEndTime := time.Now().UnixNano()
	downloadDuration := downloadEndTime - downloadStartTime

	//Hash the downloaded file
	downloadedHash := sha256.New()
	if _, err := io.Copy(downloadedHash, downloadFile); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}

	// Delete the blob
	fmt.Println("\nDeleting the blob...")
	_, err = client.DeleteBlob(context.Background(), containerName, fileName, nil)
	if err != nil {
		fmt.Println("Error deleting blob:", err)
		return
	}
	fmt.Printf("Blob '%s' deleted successfully from container '%s'.\n", fileName, containerName)

	// Delete the downloaded file
	downloadFile.Close()
	err = os.Remove(downloadedFilePath)
	if err != nil {
		fmt.Println("Error deleting downloaded file:", err)
		return
	}
	fmt.Printf("Downloaded file '%s' deleted successfully.\n", downloadedFilePath)

	// Insert results into PostgreSQL
	result := persistence.CloudResult{
		TierName:         tierName,
		Run:              run,
		Size:             fileSize,
		OriginalHash:     hex.EncodeToString(originalHash.Sum(nil)),
		DownloadedHash:   hex.EncodeToString(downloadedHash.Sum(nil)),
		StartUpload:      uploadStartTime,
		EndUpload:        uploadEndTime,
		DurationUpload:   uploadDuration,
		StartDownload:    downloadStartTime,
		EndDownload:      downloadEndTime,
		DurationDownload: downloadDuration,
	}

	if err := persistence.InsertCloudResults(result); err != nil {
		log.Fatalf("error inserting measurement: %v", err)
	}
	fmt.Printf("Inserted measurement for run %d and tier %s into PostgreSQL.\n", run, tierName)

}
