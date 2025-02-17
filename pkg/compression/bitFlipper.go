package compression

import (
	"bytes"
	"io"

	"github.com/minio/minio-go/v7"
)

type BitFlipper struct{}

func (b BitFlipper) Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, error) {
	data, err := io.ReadAll(reader)
        if err != nil {
                return nil, err
        }

        flippedData := make([]byte, len(data))
        for i, b := range data {
                flippedData[i] = ^b
        }

        return bytes.NewReader(flippedData), nil
}

func (b BitFlipper) Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, error) {
	data, err := io.ReadAll(reader)
        if err != nil {
                return nil, err
        }

        flippedData := make([]byte, len(data))
        for i, b := range data {
                flippedData[i] = ^b
        }

        return bytes.NewReader(flippedData), nil
}