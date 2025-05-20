package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type ProcessMeasurement struct {
	OriginalPath     string
	OriginalSize     int64
	OriginalHash     string

	StartGetObjectInfo    int64
	EndGetObjectInfo      int64
	DurationGetObjectInfo int64
	StartGetObject        int64
	EndGetObject          int64
	DurationGetObject     int64
	StartAlgorithm        int64
	EndAlgorithm          int64
	DurationAlgorithm     int64
	StartUpload           int64
	EndUpload             int64
	DurationUpload        int64

	Path     string
	Size     int64
	Hash     string
}

type Measurement struct {
	Algorithm string
	Run       int

	OriginalPath     string
	OriginalSize     int64
	OriginalHash     string

	Compress ProcessMeasurement

	Decompress ProcessMeasurement

	CompressionRatio float64
}

var db *sql.DB

func Connect(user string, password string, host string, port string, database string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("error opening postgres connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging postgres: %v", err)
	}
	
	return db, nil
}

func InsertMeasurement(measurement Measurement) error {
	_, err := db.Exec(`INSERT INTO measurements ( algorithm, run, original_path, original_size, original_hash, compress_start_get_object_info, compress_end_get_object_info, compress_duration_get_object_info, compress_start_get_object, compress_end_get_object, compress_duration_get_object, compress_start_algorithm, compress_end_algorithm, compress_duration_algorithm, compress_start_upload, compress_end_upload, compress_duration_upload, compress_path, compress_size, compress_hash, decompress_start_get_object_info, decompress_end_get_object_info, decompress_duration_get_object_info, decompress_start_get_object, decompress_end_get_object, decompress_duration_get_object, decompress_start_algorithm, decompress_end_algorithm, decompress_duration_algorithm, decompress_start_upload, decompress_end_upload, decompress_duration_upload, decompress_path, decompress_size, decompress_hash, compression_ratio ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36 )`,
		measurement.Algorithm, measurement.Run, measurement.OriginalPath, measurement.OriginalSize, measurement.OriginalHash, measurement.Compress.StartGetObjectInfo, measurement.Compress.EndGetObjectInfo, measurement.Compress.DurationGetObjectInfo, measurement.Compress.StartGetObject, measurement.Compress.EndGetObject, measurement.Compress.DurationGetObject, measurement.Compress.StartAlgorithm, measurement.Compress.EndAlgorithm, measurement.Compress.DurationAlgorithm, measurement.Compress.StartUpload, measurement.Compress.EndUpload, measurement.Compress.DurationUpload, measurement.Compress.Path, measurement.Compress.Size, measurement.Compress.Hash, measurement.Decompress.StartGetObjectInfo, measurement.Decompress.EndGetObjectInfo, measurement.Decompress.DurationGetObjectInfo, measurement.Decompress.StartGetObject, measurement.Decompress.EndGetObject, measurement.Decompress.DurationGetObject, measurement.Decompress.StartAlgorithm, measurement.Decompress.EndAlgorithm, measurement.Decompress.DurationAlgorithm, measurement.Decompress.StartUpload, measurement.Decompress.EndUpload, measurement.Decompress.DurationUpload, measurement.Decompress.Path, measurement.Decompress.Size, measurement.Decompress.Hash, measurement.CompressionRatio)
	return err
}

func CleanMeasurements() error {
	_, err := db.Exec(`DELETE FROM measurements`)
	return err
}