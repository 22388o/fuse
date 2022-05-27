package lightning

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

// Lightning address of the peer, in the format <pubkey>@host
type LightningAddress string

type Vertex [33]byte

type Peer struct {
	Address  string
	Inbound  bool
	PingTime time.Duration
	Pubkey   Vertex
	Sent     btcutil.Amount
	Received btcutil.Amount
}

var (
	lightningAddressRegex = `\S.*?@\S.*`
)

var (
	ErrUnknownLightningAddressFormat = errors.New("address provided does not match format <pubkey>@<host>")
)

// parseLightningAddress takes in a lightning address in format <pubkey>@<host> and parses it into its parts
func parseLightningAddress(address LightningAddress) (Vertex, string, error) {
	matched, err := regexp.MatchString(lightningAddressRegex, string(address))
	if !matched || err != nil {
		return [33]byte{}, "", ErrUnknownLightningAddressFormat
	}

	s := strings.Split(string(address), "@")
	if len(s) != 2 {
		return [33]byte{}, "", ErrUnknownLightningAddressFormat
	}

	var pubkey Vertex
	copy(pubkey[:], s[0])
	return pubkey, s[1], nil
}

// OpenChannel connects to a peer if required and opens a channel
func (l LightningClient) OpenChannel(ctx context.Context, addr LightningAddress, localSat, pushSat btcutil.Amount, private bool) (chainhash.Hash, uint32, error) {

	pubkey, host, err := parseLightningAddress(addr)
	if err != nil {
		return chainhash.Hash{}, 0, err
	}

	peers, err := l.provider.ListPeers(ctx)
	if err != nil {
		return chainhash.Hash{}, 0, err
	}

	connected := false
	for _, peer := range peers {
		if peer.Address == host && peer.Pubkey == pubkey {
			connected = true
			break
		}
	}

	if !connected {
		err := l.provider.ConnectPeer(ctx, pubkey, host)
		if err != nil {
			return chainhash.Hash{}, 0, err
		}
	}

	return l.provider.OpenChannel(ctx, pubkey, localSat, pushSat, private)
}
