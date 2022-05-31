package lightning

import (
	"context"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

// Lightning address of the peer, in the format <pubkey>@host
type LightningAddress string

const VertexSize = 33

type Vertex [VertexSize]byte

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
	ErrInvalidPubKeyLength           = errors.New("invalid pubkey length")
)

// parseLightningAddress takes in a lightning address in format <pubkey>@<host> and parses it into its parts
func parseLightningAddress(address LightningAddress) (Vertex, string, error) {
	matched, err := regexp.MatchString(lightningAddressRegex, string(address))
	if !matched || err != nil {
		return [VertexSize]byte{}, "", ErrUnknownLightningAddressFormat
	}

	s := strings.Split(string(address), "@")
	if len(s) != 2 {
		return [VertexSize]byte{}, "", ErrUnknownLightningAddressFormat
	}

	if len(s[0]) != VertexSize*2 {
		return [VertexSize]byte{}, "", ErrInvalidPubKeyLength
	}

	vertex, err := hex.DecodeString(s[0])
	if err != nil {
		return [VertexSize]byte{}, "", err
	}

	if len(vertex) != VertexSize {
		return [VertexSize]byte{}, "", ErrInvalidPubKeyLength
	}

	var pubkey Vertex
	copy(pubkey[:], vertex)
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
		if peer.Pubkey == pubkey {
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
