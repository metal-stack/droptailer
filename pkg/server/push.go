package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/metal-pod/droptailer/api/proto"
)

const (
	timestamp = "timestamp"
)

// Server is responsible to implement all droptailer interfaces
type Server struct{}

// Push implements droptailer.Push
func (s *Server) Push(ctx context.Context, de *pb.Drop) (*pb.Void, error) {
	de.Fields[timestamp] = time.Unix(de.Timestamp.Seconds, 0).String()

	js, err := json.Marshal(de.Fields)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s\n", de.Fields[timestamp], js)
	return &pb.Void{}, nil
}
