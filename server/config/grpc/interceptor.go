package grpc

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Start time for logging
	start := time.Now()
	log.Printf("🌸 gRPC Call: %s at %s", info.FullMethod, start.Format(time.RFC3339))

	// Example: Extract metadata for logging or further checks
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("🌸 Metadata received: %v", md)
	}

	// Proceed with handling the request
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("❌ Error handling request: %v", err)
	}

	// Log completion and duration
	log.Printf("🌸 Completed %s in %v with result: %v", info.FullMethod, time.Since(start), resp)
	return resp, err
}

// gracefulShutdown stops the gRPC server nicely when receiving a termination signal
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🌸 Shutting down gRPC server gracefully... ✨")
	grpcServer.GracefulStop()
	log.Println("🌸 gRPC server stopped. Bye bye~! (つ≧▽≦)つ")
}
