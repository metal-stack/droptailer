package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/metal-pod/droptailer/proto"
)

type Server struct{}

// Push implements droptailer.Push
func (s *Server) Push(ctx context.Context, de *pb.Drop) (*pb.Void, error) {
	fmt.Printf("%s %s\n", time.Unix(de.Timestamp.Seconds, 0), de.Fields)
	return &pb.Void{}, nil
}
