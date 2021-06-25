package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/metal-pod/droptailer/pkg/forwarder"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	pb "github.com/metal-pod/droptailer/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client sends drops of the journal to the droptailer server.
type Client struct {
	ServerAddress   string
	PrefixesOfDrops []string
	Certificates    Certificates
	EveSocket       string
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
	ca, err := ioutil.ReadFile(c.Certificates.CaCertificate)
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
	opts := []grpcretry.CallOption{
		grpcretry.WithBackoff(grpcretry.BackoffLinear(100 * time.Millisecond)),
	}
	conn, err := grpc.Dial(c.ServerAddress, grpc.WithTransportCredentials(creds),
		// grpc.WithStreamInterceptor(grpcretry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpcretry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return fmt.Errorf("could not connect to server: %w", err)
	}
	defer conn.Close()
	dsc := pb.NewDroptailerClient(conn)

	df, err := forwarder.NewDropforwarder(dsc, c.PrefixesOfDrops)
	if err != nil {
		return err
	}
	defer func() {
		err = df.Close()
		if err != nil {
			fmt.Printf("error closing journal reader:%s", err)
		}
	}()
	go df.Run()

	sf, err := forwarder.NewSuricataforwarder(dsc, c.EveSocket)
	if err != nil {
		return err
	}
	defer func() {
		err = sf.Close()
		if err != nil {
			fmt.Printf("error closing suricata reader:%s", err)
		}
	}()
	sf.Run()

	return nil
}
