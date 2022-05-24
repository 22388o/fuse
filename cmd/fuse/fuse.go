package main

import (
	"flag"
	"net/http"

	"github.com/btcsuite/btcutil"
	"github.com/mdedys/fuse/fuse"
	"github.com/mdedys/fuse/lightning"
	"github.com/mdedys/fuse/lnd"
	"github.com/mdedys/fuse/logging"
)

func main() {

	var (
		flags      = flag.NewFlagSet("fuse", flag.ExitOnError)
		lndAddress = flags.String("lnd-address", "localhost:10002", "Address of LND Node")
		network    = flags.String("btc-network", "regtest", "Bitcoin network to use")
		macPath    = flags.String("mac-path", "./.fuse/admin.macaroon", "Admin macaroon path")
		tlsPath    = flags.String("tls-path", "./.fuse/tls.cert", "TLS cert path")
		maxFee     = flags.Int("max-fee", 1000, "Max ln txn fee")
	)

	logging.Configure()

	client, err := lnd.NewClient(*lndAddress, *network, *macPath, *tlsPath, btcutil.Amount(*maxFee))
	if err != nil {
		panic(err)
	}

	server := fuse.New(*client, lightning.Network(*network))
	http.ListenAndServe(":1000", server)
}
