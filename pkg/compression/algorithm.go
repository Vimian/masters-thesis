package compression

import "io"

type AlgorithmFunc func(reader io.Reader, fileSize int64) (io.Reader, error)

type Algorithm interface {
	Compress() AlgorithmFunc
	Decompress() AlgorithmFunc
}