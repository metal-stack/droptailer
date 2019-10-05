package main

import (
	"log"
	"os"
	"strings"

	"github.com/metal-pod/droptailer/pkg/client"
)

const (
	defaultServerAddress = "localhost:50051"
)

var defaultPrefixesOfDrops = []string{"nftables-metal-dropped: ", "nftables-firewall-dropped: "}

func main() {
	// address should be in the form of: dns://localhost:53/droptailer:50051
	// then according to: https://github.com/grpc/grpc/blob/master/doc/naming.md
	// name based resolution should happen, which can be a /etc/hosts entry
	// which is created by the firewall-policy-controller
	// or we skip the dns resolver inbetween and just specify:
	// droptailer:50051 and rely on the local resolver which will update the IP for the client.
	// /etc/hosts can be updated with: https://github.com/txn2/txeh
	address := os.Getenv("DROPTAILER_SERVER_ADDRESS")
	if address == "" {
		address = defaultServerAddress
	}
	prefixesOfDrops := defaultPrefixesOfDrops
	prefixesOfDropsEnv := os.Getenv("DROPTAILER_PREFIXES_OF_DROPS")
	if prefixesOfDropsEnv != "" {
		prefixesOfDrops = strings.Split(prefixesOfDropsEnv, ",")
	}
	c := client.Client{
		ServerAddress:   address,
		PrefixesOfDrops: prefixesOfDrops,
	}
	err := c.Start()
	if err != nil {
		log.Fatalf("client could not start or died, %v", err)
	}
}
