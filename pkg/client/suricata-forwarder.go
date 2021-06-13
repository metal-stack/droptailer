package client

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/metal-pod/droptailer/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type suricataforwarder struct {
	dsc pb.DroptailerClient
	s   net.Conn
}

func NewSuricataforwarder(dsc pb.DroptailerClient, evesocket string) (*suricataforwarder, error) {
	c, err := net.Dial("unix", evesocket)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	df := &suricataforwarder{
		s:   c,
		dsc: dsc,
	}
	return df, nil
}
func (d *suricataforwarder) Close() error {
	return d.s.Close()
}
func (d *suricataforwarder) run() {
	d.writeTo(d.s)
}

// writeTo
// actual message will be like
func (d *suricataforwarder) writeTo(r io.ReadCloser) {
	br := bufio.NewReader(r)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			r.Close()
			break
		}
		var fields map[string]string
		err = json.Unmarshal(line, &fields)
		if err != nil {
			log.Printf("Error decoding log line to json:%v\n", err)
		}
		de := &pb.Event{
			Timestamp: timestamppb.Now(),
			Fields:    fields,
			Type:      pb.EventType_IDS,
		}
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		_, err = d.dsc.Push(
			ctx,
			de,
			grpc_retry.WithMax(30),
			grpc_retry.WithPerRetryTimeout(1*time.Second))
		if err != nil {
			log.Printf("unable to send eve entry:%v", err)
		}
	}
}
