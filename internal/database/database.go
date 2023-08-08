package database

import (
	"database/sql"
	"fmt"

	"github.com/hojamuhammet/grpc-crud-go/internal/config"
	_ "github.com/lib/pq"
)

// ConnectDatabase connects to the database using the provided configuration in config file.
// It returns a database connection object (*sql.DB) and any error encountered.
func ConnectDatabase(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
