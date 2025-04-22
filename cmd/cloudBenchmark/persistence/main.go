package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type CloudResult struct {
	TierName string
	Run      int

	Size           int64
	OriginalHash   string
	DownloadedHash string

	StartUpload    int64
	EndUpload      int64
	DurationUpload int64

	StartDownload    int64
	EndDownload      int64
	DurationDownload int64
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

func InsertCloudResults(result CloudResult) error {
	_, err := db.Exec(`INSERT INTO cloud (tier_name, run, size, original_hash, downloaded_hash, start_upload, end_upload, duration_upload, start_download, end_download, duration_download) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		result.TierName, result.Run, result.Size, result.OriginalHash, result.DownloadedHash, result.StartUpload, result.EndUpload, result.DurationUpload, result.StartDownload, result.EndDownload, result.DurationDownload)
	return err
}

func CleanCloudResults() error {
	_, err := db.Exec(`DELETE FROM cloud`)
	return err
}
