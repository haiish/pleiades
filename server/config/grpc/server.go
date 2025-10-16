package grpc

import (
	authPB "pleiades/gen/auth"
	"pleiades/server/service"

	"google.golang.org/grpc"
)

func registerServices(server *grpc.Server) {
	// pbhealth.RegisterHealthServer(server, &services.HealthService{})
	authPB.RegisterAuthServiceServer(server, &service.AuthService{})
}
