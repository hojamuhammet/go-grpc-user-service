package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/hojamuhammet/go-grpc-user-service/internal/config"
	"github.com/hojamuhammet/go-grpc-user-service/internal/database"
	"github.com/hojamuhammet/go-grpc-user-service/internal/server"
	"github.com/hojamuhammet/go-grpc-user-service/protobuf"
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

	maxPayloadSize := 1 * 1024 * 1024

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxPayloadSize),
		grpc.MaxSendMsgSize(maxPayloadSize),
	)
	protobuf.RegisterUserServiceServer(s, &server.UserServer{DB: db})

	// Enable the reflection service on the server.
	reflection.Register(s)

	log.Println("gRPC server is listening on", cfg.GRPCPort)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("Failed to serve: ", err)
		}
	}()

	// Start gRPC Gateway server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	httpMux := runtime.NewServeMux()

	err = protobuf.RegisterUserServiceHandlerFromEndpoint(context.Background(), httpMux, ":"+cfg.GRPCPort, opts)
	if err != nil {
		log.Fatal("Failed to register gRPC Gateway: ", err)
	}

	log.Println("gRPC Gateway server is listening on", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, httpMux); err != nil {
		log.Fatal("Failed to serve gRPC Gateway: ", err)
	}
}
