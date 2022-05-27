package lightning

import (
	"testing"

	"github.com/mdedys/fuse/test"
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
			addr:           "02f6725f9c1c40333b67faea92fd211c183050f28df32cac3f9d69685fe9665432@localhost:3000",
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
			test.Assert(t, tc.expectedError, err)
			test.Assert(t, tc.expectedPubkey, pubkey)
			test.Assert(t, tc.expectedHost, host)
		})
	}
}
