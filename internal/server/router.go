// server/router.go

package server

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/hojamuhammet/go-grpc-user-service/protobuf" // Import your protobuf package here
)

// RegisterHandlers registers the gRPC-Gateway handlers for your service methods.
func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC Gateway: %v", err)
	}
}

// CreateHTTPRouter creates an HTTP router with gRPC-Gateway handlers.
func CreateHTTPRouter(endpoint string, opts []grpc.DialOption) http.Handler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	RegisterHandlers(ctx, mux, endpoint, opts)

	// Optionally, you can add more HTTP routes to the mux here using standard http.HandleFunc.

	return mux
}
