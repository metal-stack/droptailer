package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/metal-pod/droptailer/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client sends drops of the journal to the droptailer server.
type Client struct {
	ServerAddress     string
	PrefixesOfDrops   []string
	PrefixesOfAccepts []string
	Certificates      Certificates
}

type Certificates struct {
	ClientCertificate string
	ClientKey         string
	CaCertificate     string
}

// Start to push drops to the droptailer server.
func (c Client) Start() error {

	// Load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(c.Certificates.ClientCertificate, c.Certificates.ClientKey)
	if err != nil {
		return fmt.Errorf("could not load client key pair: %w", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(c.Certificates.CaCertificate)
	if err != nil {
		return fmt.Errorf("could not read ca certificate: %w", err)
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
		MinVersion:   tls.VersionTLS12,
	})

	// Set up a connection to the server.
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	conn, err := grpc.NewClient(c.ServerAddress, grpc.WithTransportCredentials(creds),
		// grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return fmt.Errorf("could not connect to server: %w", err)
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
		return fmt.Errorf("error opening journal: %w", err)
	}
	if jr == nil {
		return fmt.Errorf("got a nil reader")
	}
	defer jr.Close()
	df := &dropforwarder{
		jr:             jr,
		dsc:            dsc,
		dropPrefixes:   c.PrefixesOfDrops,
		acceptPrefixes: c.PrefixesOfAccepts,
	}
	df.run()
	return nil
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
