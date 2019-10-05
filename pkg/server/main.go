package server

import (
	"context"
	"fmt"
	"time"

	"github.com/metal-pod/droptailer/droptailer"
)

type Server struct{}

// Push implements droptailer.Push
func (s *Server) Push(ctx context.Context, de *droptailer.Drop) (*droptailer.Void, error) {
	fmt.Printf("%s %s\n", time.Unix(de.Timestamp.Seconds, 0), de.Fields)
	return &droptailer.Void{}, nil
}
