package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	// Replace with your actual connection string, container name, and file path
	connectionString := "DefaultEndpointsProtocol=https;AccountName=masterthesishottier;AccountKey=GlmPvlSzv3nEUoo36RYWNsfOCfaPWrjJlg/B2oqgvzeTq/kH1RSC+oTQI9M4SnyxPGx44Bm6vJmj+AStdsXgVA==;EndpointSuffix=core.windows.net"
	containerName := "uploadcontainer"
	localFilePath := "./files/MBYNVOAVDILB2AZ6AU7FGZA3BWG3XXZ4.pdf"
	blobName := "MBYNVOAVDILB2AZ6AU7FGZA3BWG3XXZ4.pdf"

	// Create a BlobServiceClient
	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		fmt.Println("Error creating service client:", err)
		return
	}

	// Get a BlobContainerClient
	containerClient := serviceClient.ServiceClient().NewContainerClient(containerName)

	// Get a BlockBlobClient
	blobClient := containerClient.NewBlockBlobClient(blobName)

	// *** Upload the file ***
	fmt.Println("Starting file upload...")

	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return
	}
	defer file.Close()

	uploadStartTime := time.Now()

	_, err = blobClient.UploadStream(context.Background(), file, nil)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	uploadEndTime := time.Now()
	uploadDuration := uploadEndTime.Sub(uploadStartTime)
	fmt.Printf("File uploaded successfully to '%s' in container '%s'.\n", blobName, containerName)
	fmt.Printf("Upload started at: %s\n", uploadStartTime.Format(time.RFC3339))
	fmt.Printf("Upload finished at: %s\n", uploadEndTime.Format(time.RFC3339))
	fmt.Printf("Upload duration: %s\n", uploadDuration)

	// *** Download the file ***
	fmt.Println("\nStarting file download...")

	downloadedFilePath := "./downloads/downloaded_" + blobName
	downloadFile, err := os.Create(downloadedFilePath)
	if err != nil {
		fmt.Println("Error creating download file:", err)
		return
	}
	defer downloadFile.Close()

	downloadStartTime := time.Now()

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

	downloadEndTime := time.Now()
	downloadDuration := downloadEndTime.Sub(downloadStartTime)
	fmt.Printf("File downloaded successfully to '%s'.\n", downloadedFilePath)
	fmt.Printf("Download started at: %s\n", downloadStartTime.Format(time.RFC3339))
	fmt.Printf("Download finished at: %s\n", downloadEndTime.Format(time.RFC3339))
	fmt.Printf("Download duration: %s\n", downloadDuration)

	// *** Delete the blob ***
	fmt.Println("\nDeleting the blob...")
	_, err = blobClient.Delete(context.Background(), nil)
	if err != nil {
		fmt.Println("Error deleting blob:", err)
		return
	}
	fmt.Printf("Blob '%s' deleted successfully from container '%s'.\n", blobName, containerName)

	downloadFile.Close()
	err = os.Remove(downloadedFilePath)
	if err != nil {
		fmt.Println("Error deleting downloaded file:", err)
		return
	}
}
