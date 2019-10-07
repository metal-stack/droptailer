package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/metal-pod/droptailer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client sends drops of the journal to the droptailer server.
type Client struct {
	ServerAddress   string
	PrefixesOfDrops []string
}

// Start to push drops to the droptailer server.
func (c Client) Start() error {

	// Load the certificates from disk
	certificate, err := tls.LoadX509KeyPair("./cfssl/client.pem", "./cfssl/client-key.pem")
	if err != nil {
		return fmt.Errorf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./cfssl/ca.pem")
	if err != nil {
		return fmt.Errorf("could not read ca certificate: %s", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return errors.New("failed to append ca certs")
	}

	// Create the TLS credentials for transport
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   "droptailer",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// Set up a connection to the server.
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	conn, err := grpc.Dial(c.ServerAddress, grpc.WithTransportCredentials(creds),
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
	msg, ok := entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]
	if !ok {
		return "", fmt.Errorf("no %s field present in journal entry", sdjournal.SD_JOURNAL_FIELD_MESSAGE)
	}
	usec := entry.RealtimeTimestamp
	timestamp := time.Unix(0, int64(usec)*int64(time.Microsecond))
	return fmt.Sprintf("%d@%s\n", timestamp.Unix(), msg), nil
}
