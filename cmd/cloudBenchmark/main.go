package main

import (
	"bytes"
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
	"github.com/vimian/masters-thesis/cmd/cloudbenchmark/persistence"
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

	//test(numberOfRuns, "premium", "DefaultEndpointsProtocol=https;AccountName=masterthesispremiumtier;AccountKey=iqQX4EujDpxb2+L64cAsEhy96aduGbgWh/AdcweoZFsnwVZS8P0Wmr/QPBiqadaq9rsn2xac+05P+AStV0atxg==;EndpointSuffix=core.windows.net", filesnames)
	test(numberOfRuns, "hot", "DefaultEndpointsProtocol=https;AccountName=masterthesishottier;AccountKey=GlmPvlSzv3nEUoo36RYWNsfOCfaPWrjJlg/B2oqgvzeTq/kH1RSC+oTQI9M4SnyxPGx44Bm6vJmj+AStdsXgVA==;EndpointSuffix=core.windows.net", filesnames)
	//test(numberOfRuns, "cool", "DefaultEndpointsProtocol=https;AccountName=masterthesiscooltier;AccountKey=+RF+VOPKw0Ju7L3n/POuA4zbLWYwMYud9CWQEgcQ9ZheBQbN2Nnx0CDRmDZmpTVx8OlBrXlmuAe8+AStYw9tqQ==;EndpointSuffix=core.windows.net", filesnames)
	//test(numberOfRuns, "cold", "DefaultEndpointsProtocol=https;AccountName=masterthesiscoldtier;AccountKey=1lyP4P4LQzbsV/PB0jsFyO77Gm/CVnQvDAfOF87en4J2U8+dGSyG0aAfebvMAqX/wKXIsxRckX/k+AStn8nDNA==;EndpointSuffix=core.windows.net", filesnames)

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

func test(numberOfRuns int, tierName string, connectionString string, fileNames []string) {
	for i := 0; i < numberOfRuns; i++ {
		for _, fileName := range fileNames {
			run(i, tierName, connectionString, fileName)
		}
	}
}

func run(run int, tierName string, connectionString string, fileName string) {
	containerName := "uploadcontainer"
	localFilePath := "./files/" + fileName

	// Create a BlobServiceClient
	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		fmt.Println("Error creating service client:", err)
		return
	}

	// Get BlobContainerClient
	containerClient := serviceClient.ServiceClient().NewContainerClient(containerName)

	// Get BlockBlobClient
	blobClient := containerClient.NewBlockBlobClient(fileName)

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
	originalHash := sha256.New()
	if _, err := io.Copy(originalHash, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}
	fmt.Println("original hash:", hex.EncodeToString(originalHash.Sum(nil)))
	uploadStartTime := time.Now().UnixNano()

	_, err = blobClient.UploadStream(context.Background(), file, nil)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	uploadEndTime := time.Now().UnixNano()
	uploadDuration := uploadEndTime - uploadStartTime
	fmt.Printf("File uploaded successfully to '%s' in container '%s'.\n", fileName, containerName)

	// Download the file
	fmt.Println("\nStarting file download...")

	downloadedFilePath := "./downloads/" + fileName
	downloadFile, err := os.Create(downloadedFilePath)
	if err != nil {
		fmt.Println("Error creating download file:", err)
		return
	}
	defer downloadFile.Close()

	downloadStartTime := time.Now().UnixNano()
	/*
		get, err := blobClient.DownloadStream(context.Background(), nil)
		if err != nil {
			fmt.Println("Error downloading file:", err)
			return
		}
		defer get.Body.Close()

		_, err = io.Copy(downloadFile, get.Body)
		if err != nil {
			fmt.Println("Error copying downloaded data to file:", err)
			return
		}
	*/

	// Download the blob
	//TODO
	//LOOK HERE
	get, err := blobClient.DownloadStream(context.Background(), nil)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(context.Background(), &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		fmt.Println("Error reading downloaded data:", err)
		return
	}

	err = retryReader.Close()
	if err != nil {
		fmt.Println("Error closing retry reader:", err)
		return
	}

	//Hash the downloaded file
	downloadedHash := sha256.New()
	if _, err := io.Copy(downloadedHash, &downloadedData); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}

	downloadEndTime := time.Now().UnixNano()
	downloadDuration := downloadEndTime - downloadStartTime
	fmt.Printf("File downloaded successfully to '%s'.\n", downloadedFilePath)

	// Delete the blob
	fmt.Println("\nDeleting the blob...")
	_, err = blobClient.Delete(context.Background(), nil)
	if err != nil {
		fmt.Println("Error deleting blob:", err)
		return
	}
	fmt.Printf("Blob '%s' deleted successfully from container '%s'.\n", fileName, containerName)

	// Delete the downloaded file
	/*
		downloadFile.Close()
		err = os.Remove(downloadedFilePath)
		if err != nil {
			fmt.Println("Error deleting downloaded file:", err)
			return
		}
		fmt.Printf("Downloaded file '%s' deleted successfully.\n", downloadedFilePath)
	*/

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
