package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type AnalysisResult struct {
	FilePath string
	FileName string
	FileSize int64
	WindowLengthBytes int64
	DictionaryLength int64
	DictionaryLimitReached int8
	CompressedSizeBytes int64
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

func InsertAnalysisResult(result AnalysisResult) error {
	_, err := db.Exec(`INSERT INTO analytics ( file_path, file_name, file_size, window_length_bytes, dictionary_length, dictionary_limit_reached, compressed_size_bytes_go ) VALUES ( $1, $2, $3, $4, $5, $6, $7 )`,
		result.FilePath, result.FileName, result.FileSize, result.WindowLengthBytes, result.DictionaryLength, result.DictionaryLimitReached, result.CompressedSizeBytes)
	return err
}

func CleanAnalytics() error {
	_, err := db.Exec(`DELETE FROM analytics`)
	return err
}