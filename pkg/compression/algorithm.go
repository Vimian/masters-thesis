package compression

import (
	"io"

	"github.com/minio/minio-go/v7"
)

type AlgorithmFunc func(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error)

type Algorithm interface {
	Compress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error)
	Decompress(reader io.Reader, objectInfo minio.ObjectInfo) (io.Reader, string, error)
}