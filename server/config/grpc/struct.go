package grpc

import (
	"net"

	"google.golang.org/grpc"
)

var grpcServer *grpc.Server

type ipFilteredListener struct {
	net.Listener
}
