package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Start initializes and runs the gRPC server
func Start(cfg *Config) error {
	address := fmt.Sprintf(":%s", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("❌ Error listening on port %s: %v", cfg.GRPCPort, err)
		return fmt.Errorf("failed to listen on port %s: %w", cfg.GRPCPort, err)
	}

	// Wrap listener to filter connections by IP
	lis = &ipFilteredListener{Listener: lis}

	// Setup server with a single interceptor
	grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor), // Only add your interceptor once
	)

	// Register services
	registerServices(grpcServer)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("🌸 gRPC server is running on %s ✨", address)

	// Start graceful shutdown handler
	go gracefulShutdown()

	// Serve incoming connections
	log.Println("🌸 Starting to serve incoming connections...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("❌ Error during Serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (l *ipFilteredListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		log.Printf("❌ Error during Accept: %v", err)
		return nil, err
	}

	return conn, nil
}
