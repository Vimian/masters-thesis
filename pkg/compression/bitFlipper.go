package compression

import (
	"bytes"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

type BitFlipper struct{}

func (b BitFlipper) Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	data, err := io.ReadAll(reader)
        if err != nil {
                return nil, "", err
        }

        flippedData := make([]byte, len(data))
        for i, b := range data {
                flippedData[i] = ^b
        }

        var parts []string = strings.Split(objectInfo.Key, "/")
	var fileName string = parts[len(parts)-1]

        return bytes.NewReader(flippedData), fileName, nil
}

func (b BitFlipper) Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error) {
	data, err := io.ReadAll(reader)
        if err != nil {
                return nil, "", err
        }

        flippedData := make([]byte, len(data))
        for i, b := range data {
                flippedData[i] = ^b
        }

        var parts []string = strings.Split(objectInfo.Key, "/")
	var fileName string = parts[len(parts)-1]

        return bytes.NewReader(flippedData), fileName, nil
}