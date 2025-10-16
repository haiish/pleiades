package grpc

import "os"

type Config struct {
	GRPCPort string
}

func LoadConfig() *Config {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "3552" // default port if not set, safety net! 🎀
	}

	return &Config{
		GRPCPort: port,
	}
}
