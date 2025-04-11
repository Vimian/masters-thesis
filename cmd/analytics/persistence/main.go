package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type AnalysisResult struct {
	FilePath string
	Bytes     int64
	BytesNeeded int64
	DictionarySize int64
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
	_, err := db.Exec(`INSERT INTO analytics ( file_path, bytes, bytes_needed, dictionary_size ) VALUES ( $1, $2, $3, $4 )`,
		result.FilePath, result.Bytes, result.BytesNeeded, result.DictionarySize)
	return err
}

func CleanAnalytics() error {
	_, err := db.Exec(`DELETE FROM analytics`)
	return err
}