package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	pb "github.com/metal-pod/droptailer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	server "github.com/metal-pod/droptailer/pkg/server"
)

const (
	defaultServerCertificate = "tls.crt"
	defaultServerKey         = "tls.key"
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

	// Read cert and key file
	serverCertificate := os.Getenv("SERVER_CERTIFICATE")
	if serverCertificate == "" {
		serverCertificate = defaultServerCertificate
	}
	serverKey := os.Getenv("SERVER_KEY")
	if serverKey == "" {
		serverKey = defaultServerKey
	}

	// Generate Certificate struct
	c, err := tls.LoadX509KeyPair(serverCertificate, serverKey)
	if err != nil {
		log.Fatalf("failed to parse certificate: %v", err)
	}
	// Create the TLS credentials
	creds := credentials.NewServerTLSFromCert(&c)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}
	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	pb.RegisterDroptailerServer(s, &server.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
