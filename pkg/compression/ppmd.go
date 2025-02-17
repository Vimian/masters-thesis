package compression

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/minio/minio-go/v7"
)

type LZMA struct{}

func (l LZMA) Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fmt.Println("compressing", objectInfo.Key) // TODO: remove this
	if err = os.WriteFile(objectInfo.Key, data, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, err
	}

	cmd := exec.Command("7zip", "a", "-t7z", objectInfo.Key + ".7z", objectInfo.Key)
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	compressedData, err := os.ReadFile(objectInfo.Key + ".7z")

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(compressedData), nil
}

func (l LZMA) Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, error) {
	compressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fmt.Println("decompressing", objectInfo.Key) // TODO: remove this
	if err = os.WriteFile(objectInfo.Key, compressedData, 0644); err != nil { // TODO: maybe 0777 permissions is needed...
		return nil, err
	}

	cmd := exec.Command("7zip", "e", objectInfo.Key)
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	//TODO: find the file in the directory vv maybe use the objectInfo.Key and find the file witt other name
	data, err := os.ReadFile(objectInfo.Key)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}