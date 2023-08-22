// server/router.go

package server

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(customHTTPErrorMapper),
	)
	RegisterHandlers(ctx, mux, endpoint, opts)

	return mux
}


func customHTTPErrorMapper(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
    // Handle gRPC errors and map them to HTTP status codes.
    grpcStatus := status.Convert(err)
    switch grpcStatus.Code() {
    case codes.NotFound:
        w.WriteHeader(http.StatusNotFound)
    default:
        w.WriteHeader(http.StatusInternalServerError)
    }
}