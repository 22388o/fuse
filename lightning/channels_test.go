package lightning

import (
	"context"
	"errors"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/mdedys/fuse/test"
)

var (
	addr = LightningAddress("02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432@localhost:3000")
)

func TestParseLightningAddress(t *testing.T) {
	type testCase struct {
		name           string
		addr           LightningAddress
		expectedPubkey Vertex
		expectedHost   string
		expectedError  error
	}

	var validPubkey Vertex
	copy(validPubkey[:], "02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432")

	tests := []testCase{
		{
			name:           "valid address",
			addr:           addr,
			expectedPubkey: validPubkey,
			expectedHost:   "localhost:3000",
			expectedError:  nil,
		},
		{
			name:           "invalid address no @",
			addr:           "02f6725f9c1c40333",
			expectedPubkey: [33]byte{},
			expectedHost:   "",
			expectedError:  ErrUnknownLightningAddressFormat,
		},
		{
			name:           "invalid address no host",
			addr:           "02f6725f9c1c40333@",
			expectedPubkey: [33]byte{},
			expectedHost:   "",
			expectedError:  ErrUnknownLightningAddressFormat,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pubkey, host, err := parseLightningAddress(tc.addr)
			test.AssertEqual(t, tc.expectedError, err)
			test.AssertEqual(t, tc.expectedPubkey, pubkey)
			test.AssertEqual(t, tc.expectedHost, host)
		})
	}
}

func TestOpenChann(t *testing.T) {
	type testCase struct {
		name         string
		addr         LightningAddress
		localSats    btcutil.Amount
		pushSat      btcutil.Amount
		expectedHash chainhash.Hash
		expectError  bool
		mocks        mockLightningProvider
	}

	validChainHash, _ := chainhash.NewHashFromStr("14a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef")

	tests := []testCase{
		{
			name:      "happy path no peers",
			addr:      addr,
			localSats: 1000,
			pushSat:   500,
			mocks: mockLightningProvider{
				listPeers: func(ctx context.Context) ([]Peer, error) {
					return []Peer{}, nil
				},
				connectPeer: func(ctx context.Context, peer Vertex, host string) error {
					return nil
				},
				openChannel: func(ctx context.Context, peer Vertex, localSat, pushSat btcutil.Amount, private bool) (chainhash.Hash, uint32, error) {
					return *validChainHash, 1, nil
				},
			},
			expectedHash: *validChainHash,
			expectError:  false,
		},
		{
			name:      "failed to list peers",
			addr:      addr,
			localSats: 1000,
			pushSat:   500,
			mocks: mockLightningProvider{
				listPeers: func(ctx context.Context) ([]Peer, error) {
					return []Peer{}, errors.New("explosion")
				},
			},
			expectedHash: chainhash.Hash{},
			expectError:  true,
		},
		{
			name:      "failed to connect to peer",
			addr:      addr,
			localSats: 1000,
			pushSat:   500,
			mocks: mockLightningProvider{
				listPeers: func(ctx context.Context) ([]Peer, error) {
					return []Peer{}, nil
				},
				connectPeer: func(ctx context.Context, peer Vertex, host string) error {
					return errors.New("explosion")
				},
			},
			expectedHash: chainhash.Hash{},
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			client := New(tc.mocks)
			hash, _, err := client.OpenChannel(ctx, tc.addr, tc.localSats, tc.pushSat, false)

			test.AssertEqual(t, tc.expectedHash, hash)

			if !tc.expectError {
				test.AssertNil(t, err)
			} else {
				test.AssertDefined(t, err)
			}
		})
	}
}
