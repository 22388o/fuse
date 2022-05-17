package main

import (
	"flag"
	"net/http"

	"github.com/mdedys/fuse/fuse"
	"github.com/mdedys/fuse/lnd"
)

func main() {

	var (
		flags      = flag.NewFlagSet("fuse", flag.ExitOnError)
		lndAddress = flags.String("lnd-address", "", "Address of LND Node")
		network    = flags.String("btc-network", "regtest", "Bitcoin network to use")
		macPath    = flags.String("mac-path", "", "Admin macaroon path")
		tlsPath    = flags.String("tls-path", "", "TLS cert path")
	)

	client, err := lnd.NewClient(*lndAddress, *network, *macPath, *tlsPath)
	if err != nil {
		panic(err)
	}

	server := fuse.New(*client)
	http.ListenAndServe(":1000", server)
}
