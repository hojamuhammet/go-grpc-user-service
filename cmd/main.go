package main

import (
	"log"
	"net"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hojamuhammet/grpc-crud-go/internal/config"
	"github.com/hojamuhammet/grpc-crud-go/internal/database"
	"github.com/hojamuhammet/grpc-crud-go/internal/server"
	pb "github.com/hojamuhammet/grpc-crud-go/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Read the environment variables from .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Init the database connection
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Testing connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}
	log.Println("Connected to the database")

	// Start the gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server.UserServer{DB: db})

	// Enable the reflection service on the server.
	reflection.Register(s)

	log.Println("gRPC server is listening on", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
