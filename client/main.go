package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/metal-pod/droptailer/pkg/client"
)

const (
	defaultServerAddress     = "localhost:50051"
	defaultCaCertificate     = "/etc/droptailer-client/ca.crt"
	defaultClientCertificate = "/etc/droptailer-client/tls.crt"
	defaultClientKey         = "/etc/droptailer-client/tls.key"
)

var defaultPrefixesOfDrops = []string{"nftables-metal-dropped: ", "nftables-firewall-dropped: "}
var defaultPrefixesOfAccepts = []string{"nftables-metal-accept: ", "nftables-firewall-accept: "}

func main() {
	// address should be in the form of: dns://localhost:53/droptailer:50051
	// then according to: https://github.com/grpc/grpc/blob/master/doc/naming.md
	// name based resolution should happen, which can be a /etc/hosts entry
	// which is created by the firewall-controller
	// or we skip the dns resolver inbetween and just specify:
	// droptailer:50051 and rely on the local resolver which will update the IP for the client.
	// /etc/hosts can be updated with: https://github.com/txn2/txeh
	address := os.Getenv("DROPTAILER_SERVER_ADDRESS")
	if address == "" {
		address = defaultServerAddress
	}

	_, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		log.Fatalf("could not reach droptailer server within 5 seconds: %v", err)
	}

	prefixesOfDrops := defaultPrefixesOfDrops
	prefixesOfDropsEnv := os.Getenv("DROPTAILER_PREFIXES_OF_DROPS")
	if prefixesOfDropsEnv != "" {
		prefixesOfDrops = strings.Split(prefixesOfDropsEnv, ",")
	}

	prefixesOfAccepts := defaultPrefixesOfAccepts
	prefixesOfAcceptsEnv := os.Getenv("DROPTAILER_PREFIXES_OF_ACCEPTS")
	if prefixesOfAcceptsEnv != "" {
		prefixesOfAccepts = strings.Split(prefixesOfAcceptsEnv, ",")
	}

	caCertificate := os.Getenv("DROPTAILER_CA_CERTIFICATE")
	if caCertificate == "" {
		caCertificate = defaultCaCertificate
	}
	clientCertificate := os.Getenv("DROPTAILER_CLIENT_CERTIFICATE")
	if clientCertificate == "" {
		clientCertificate = defaultClientCertificate
	}
	clientKey := os.Getenv("DROPTAILER_CLIENT_KEY")
	if clientKey == "" {
		clientKey = defaultClientKey
	}
	c := client.Client{
		ServerAddress:     address,
		PrefixesOfDrops:   prefixesOfDrops,
		PrefixesOfAccepts: prefixesOfAccepts,
		Certificates: client.Certificates{
			CaCertificate:     caCertificate,
			ClientCertificate: clientCertificate,
			ClientKey:         clientKey,
		},
	}
	err = c.Start()
	if err != nil {
		log.Fatalf("client could not start or died, %v", err)
	}
}
