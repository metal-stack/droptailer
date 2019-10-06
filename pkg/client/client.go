package client

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/metal-pod/droptailer/proto"

	"google.golang.org/grpc"
)

// Client sends drops of the journal to the droptailer server.
type Client struct {
	ServerAddress   string
	PrefixesOfDrops []string
}

// Start to push drops to the droptailer server.
func (c Client) Start() error {
	// Set up a connection to the server.
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	conn, err := grpc.Dial(c.ServerAddress, grpc.WithInsecure(),
		// grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}
	defer conn.Close()
	dsc := pb.NewDroptailerClient(conn)
	jr, err := sdjournal.NewJournalReader(
		sdjournal.JournalReaderConfig{
			NumFromTail: 100,
			// Matches on message only match the whole message not the start
			Matches: []sdjournal.Match{
				{
					Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
					Value: "kernel",
				},
			},
			Formatter: messageFormatter,
		})
	if err != nil {
		return fmt.Errorf("Error opening journal: %s", err)
	}
	if jr == nil {
		return fmt.Errorf("Got a nil reader")
	}
	defer jr.Close()
	df := &dropforwarder{
		jr:       jr,
		dsc:      dsc,
		prefixes: c.PrefixesOfDrops,
	}
	df.run()
	return nil
}

func newCtx(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), timeout)
	return ctx
}

func messageFormatter(entry *sdjournal.JournalEntry) (string, error) {
	msg, ok := entry.Fields["MESSAGE"]
	if !ok {
		return "", fmt.Errorf("no MESSAGE field present in journal entry")
	}
	usec := entry.RealtimeTimestamp
	timestamp := time.Unix(0, int64(usec)*int64(time.Microsecond))
	return fmt.Sprintf("%d@%s\n", timestamp.Unix(), msg), nil
}
