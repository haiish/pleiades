package service

import authGRPC "pleiades/gen/auth"

type AuthService struct {
	authGRPC.UnimplementedAuthServiceServer
}
