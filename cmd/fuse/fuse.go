package main

import (
	"flag"
	"net/http"

	"github.com/btcsuite/btcutil"
	"github.com/mdedys/fuse/fuse"
	"github.com/mdedys/fuse/lightning"
	"github.com/mdedys/fuse/lnd"
	"github.com/mdedys/fuse/lnurl"
	"github.com/mdedys/fuse/logging"
	"github.com/rs/zerolog/log"
)

func main() {

	var (
		flags      = flag.NewFlagSet("fuse", flag.ExitOnError)
		lndAddress = flags.String("lnd-address", "localhost:1000", "Address of LND Node")
		network    = flags.String("btc-network", "regtest", "Bitcoin network to use")
		macPath    = flags.String("mac-path", "/Users/mikededys/github/fuse/.fuse/admin.macaroon", "Admin macaroon path")
		tlsPath    = flags.String("tls-path", "/Users/mikededys/github/fuse/.fuse/tls.cert", "TLS cert path")
		maxFee     = flags.Int("max-fee", 1000, "Max ln txn fee")
	)

	logging.Configure()

	log.Info().Msg("Starting fuse setup")

	client, err := lnd.NewClient(*lndAddress, *network, *macPath, *tlsPath, btcutil.Amount(*maxFee))
	if err != nil {
		panic(err)
	}

	store := lnurl.NewStore()

	log.Info().Msg("Lightning Node connection successful")

	provider := lightning.New(*client)
	server := fuse.New(provider, lightning.Network(*network), store)

	log.Info().Msg("Starting server on PORT: 1100")
	http.ListenAndServe(":1100", server)
}
