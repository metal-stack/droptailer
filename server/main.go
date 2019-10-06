package main

import (
	"log"
	"net"
	"os"

	pb "github.com/metal-pod/droptailer/proto"

	server "github.com/metal-pod/droptailer/pkg/server"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s\n", port)
	s := grpc.NewServer()
	pb.RegisterDroptailerServer(s, &server.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
