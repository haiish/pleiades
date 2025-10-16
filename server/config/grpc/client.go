package grpc

import (
	"log"
	"math"
	"pleiades/gen/auth"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	AuthService auth.AuthServiceClient
	authConn    *grpc.ClientConn
)

// Connect initializes the AuthService client and starts a background retry loop.
// It logs connection status but never crashes the app.
func Connect() error {
	const addr = "localhost:3551"

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("⚠️ Failed to create gRPC client for AuthService (%s): %v", addr, err)
		return err
	}

	authConn = conn
	AuthService = auth.NewAuthServiceClient(conn)

	// Start background connection + retry routine
	go maintainAuthConnection(addr, conn)

	return nil
}

// maintainAuthConnection periodically checks and retries connection if needed.
func maintainAuthConnection(addr string, conn *grpc.ClientConn) {
	var retryCount int

	for {
		state := conn.GetState()

		if state == connectivity.Ready {
			log.Printf("✅ AuthService connected at %s", addr)
			time.Sleep(10 * time.Second) // check periodically
			continue
		}

		retryCount++
		backoff := time.Duration(math.Min(float64(retryCount*2), 30)) * time.Second // exponential up to 30s
		log.Printf("🟡 AuthService not ready (state: %s). Retrying in %v...", state, backoff)

		conn.Connect() // non-blocking trigger
		time.Sleep(backoff)
	}
}

// IsAuthServiceAvailable checks if AuthService client has been initialized.
func IsAuthServiceAvailable() bool {
	return AuthService != nil && authConn != nil
}

// CloseAuthConnection closes the AuthService connection safely.
func CloseAuthConnection() {
	if authConn != nil {
		_ = authConn.Close()
	}
}
