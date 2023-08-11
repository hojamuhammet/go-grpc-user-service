package database

import (
	"testing"

	"github.com/hojamuhammet/go-grpc-user-service/internal/config"
)

func TestConnectDatabase(t *testing.T) {
    // Create a mock configuration
    cfg := &config.Config{
        DBHost:     "localhost",
        DBPort:     "5432",
        DBUser:     "postgres",
        DBPassword: "K862008971a!",
        DBName:     "user_administration",
    }

    // Call the function to be tested
    db, err := ConnectDatabase(cfg)

    // Check for errors
    if err != nil {
        t.Fatalf("Error connecting to database: %v", err)
    }
    defer db.Close() // Close the connection when done

    // Check if the db connection is not nil
    if db == nil {
        t.Fatal("Database connection is nil")
    }
}
