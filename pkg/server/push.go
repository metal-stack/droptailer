package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/metal-pod/droptailer/proto"
)

type Server struct{}

// Push implements droptailer.Push
func (s *Server) Push(ctx context.Context, de *pb.Drop) (*pb.Void, error) {
	de.Fields["timestamp"] = fmt.Sprintf("%s", time.Unix(de.Timestamp.Seconds, 0))

	js, err := json.Marshal(de.Fields)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s\n", de.Fields["timestamp"], js)
	return &pb.Void{}, nil
}
