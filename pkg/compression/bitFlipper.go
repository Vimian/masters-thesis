package compression

import (
	"bytes"
	"io"
)

type BitFlipper struct{}

func (b BitFlipper) Compress(reader io.Reader, fileSize int64) (io.Reader, error) {
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

func (b BitFlipper) Decompress(reader io.Reader, fileSize int64) (io.Reader, error) {
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