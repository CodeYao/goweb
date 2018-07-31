package ca_grpcserver

import (
	pb "ca/goweb/ca_grpc"
	"ca/goweb/models"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":9092"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) JudgeAddress(ctx context.Context, in *pb.AddressRequest) (*pb.AddressReply, error) {
	return &pb.AddressReply{Message: "Hello " + in.Addr}, nil
}

func (s *server) AuthorityCtrl(ctx context.Context, in *pb.AddressRequest) (*pb.IsPermissionReply, error) {
	addresslist := models.QueryData("address")
	for _, v := range addresslist {
		if in.Addr == v["address"] {
			fmt.Println("chenyao**************true")
			return &pb.IsPermissionReply{IsPermission: true}, nil
		}
	}
	fmt.Println("chenyao**************false")
	return &pb.IsPermissionReply{IsPermission: false}, nil
}

func CAGrpcRun() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}
