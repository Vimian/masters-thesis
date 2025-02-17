package compression

import (
	"io"

	"github.com/minio/minio-go/v7"
)

type AlgorithmFunc func(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, error)

type Algorithm interface {
	Compress() AlgorithmFunc
	Decompress() AlgorithmFunc
}