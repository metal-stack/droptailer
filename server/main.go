/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../journal --go_out=plugins=grpc:../journal ../journal/journal.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/metal-pod/droptailer/dropsink"
	pb "github.com/metal-pod/droptailer/dropsink"
	"google.golang.org/grpc"
)

// server is used to implement dropsink.Push
type server struct{}

// Push implements dropsink.Push
func (s *server) Push(ctx context.Context, de *dropsink.DropEntry) (*dropsink.Void, error) {
	fmt.Printf("%s %s\n", time.Unix(de.Timestamp, 0), de.Fields)
	return &dropsink.Void{}, nil
}

func main() {
	port := os.Getenv("SERVER_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s\n", port)
	s := grpc.NewServer()
	pb.RegisterDropSinkServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
