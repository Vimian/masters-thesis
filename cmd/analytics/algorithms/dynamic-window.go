package algorithms

import (
	"github.com/minio/minio-go/v7"
	"github.com/vimian/masters-thesis/cmd/analytics/persistence"
)

type DynamicWindow struct {
}

func (d DynamicWindow) AnalyzeFile(reader *minio.Object, filePath string, fileInfo minio.ObjectInfo) ([]persistence.AnalysisResult, error) {
	return nil, nil
}