package algorithms

import (
	"github.com/minio/minio-go/v7"
	"github.com/vimian/masters-thesis/cmd/analytics/persistence"
)

type Algorithm interface {
	AnalyzeFile(reader *minio.Object, filePath string, fileInfo minio.ObjectInfo) ([]persistence.AnalysisResult, error)
}