package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/metal-pod/droptailer/proto"
)

const (
	timestamp = "timestamp"
)

// Server is responsible to implement all droptailer interfaces
type Server struct{}

// Push implements droptailer.Push
func (s *Server) Push(ctx context.Context, e *pb.Event) (*pb.Void, error) {
	e.Fields[timestamp] = time.Unix(e.Timestamp.Seconds, 0).String()

	js, err := json.Marshal(e.Fields)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s %s\n", e.Type, e.Fields[timestamp], js)
	return &pb.Void{}, nil
}
