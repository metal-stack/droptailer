package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"os"

	pb "github.com/metal-pod/droptailer/proto"

	server "github.com/metal-pod/droptailer/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	serverCert, _ := ioutil.ReadFile("./cfssl/server.pem")
	serverKey, _ := ioutil.ReadFile("./cfssl/server-key.pem")

	// Generate Certificate struct
	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("failed to parse certificate: %v", err)
	}
	// // Create the TLS credentials
	creds := credentials.NewServerTLSFromCert(&cert)
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
