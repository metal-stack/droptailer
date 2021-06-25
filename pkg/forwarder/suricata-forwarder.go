package forwarder

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/metal-pod/droptailer/proto"
)

const (
	udpPacketPayloadMaxSize = 65507
)

type suricataforwarder struct {
	dsc  pb.DroptailerClient
	conn net.PacketConn
}

func NewSuricataforwarder(dsc pb.DroptailerClient, evesocket string) (*suricataforwarder, error) {
	// Delete old socket file if exists to avoid problems with ownership
	os.Remove(evesocket)

	c, err := net.ListenPacket("unixgram", evesocket)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	df := &suricataforwarder{
		dsc:  dsc,
		conn: c,
	}
	return df, nil
}
func (d *suricataforwarder) Close() error {
	return d.conn.Close()
}
func (d *suricataforwarder) Run() {
	defer d.Close()

	for {
		payload := make([]byte, udpPacketPayloadMaxSize)
		n, _, err := d.conn.ReadFrom(payload)
		if n < 1 {
			continue
		}
		if err != nil {
			fmt.Printf("Error while reading from suricata output: %s\n", err)
			break
		}

		de := &pb.Event{
			Timestamp: timestamppb.Now(),
			Type:      pb.EventType_IDS,
			Content:   string(payload[:n]),
		}
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		defer cancel()
		if _, err := d.dsc.Push(
			ctx,
			de,
			grpcretry.WithMax(30),
			grpcretry.WithPerRetryTimeout(1*time.Second),
		); err != nil {
			log.Printf("unable to send eve entry:%v", err)
		}
	}
}
