package lightning

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
)

type mockLightningProvider struct {
	walletBalance func(ctx context.Context) (btcutil.Amount, error)

	addInvoice func(ctx context.Context, value lnwire.MilliSatoshi, memo string) (Invoice, error)
	payInvoice func(ctx context.Context, invoice Invoice) (PaymentResult, error)

	listPeers   func(ctx context.Context) ([]Peer, error)
	connectPeer func(ctx context.Context, peer Vertex, host string) error
	openChannel func(ctx context.Context, peer Vertex, localSat, pushSat btcutil.Amount, private bool) (chainhash.Hash, uint32, error)
}

func (m mockLightningProvider) WalletBalance(ctx context.Context) (btcutil.Amount, error) {
	if m.walletBalance != nil {
		return m.walletBalance(ctx)
	}
	return btcutil.Amount(0), nil
}

func (m mockLightningProvider) AddInvoice(ctx context.Context, value lnwire.MilliSatoshi, memo string) (Invoice, error) {
	if m.addInvoice != nil {
		return m.addInvoice(ctx, value, memo)
	}
	return Invoice{}, nil
}

func (m mockLightningProvider) PayInvoice(ctx context.Context, invoice Invoice) (PaymentResult, error) {
	if m.payInvoice != nil {
		return m.payInvoice(ctx, invoice)
	}
	return PaymentResult{}, nil
}

func (m mockLightningProvider) ListPeers(ctx context.Context) ([]Peer, error) {
	if m.listPeers != nil {
		return m.listPeers(ctx)
	}
	return []Peer{}, nil
}

func (m mockLightningProvider) ConnectPeer(ctx context.Context, peer Vertex, host string) error {
	if m.connectPeer != nil {
		return m.connectPeer(ctx, peer, host)
	}
	return nil
}

func (m mockLightningProvider) OpenChannel(ctx context.Context, peer Vertex, localSat, pushSat btcutil.Amount, private bool) (chainhash.Hash, uint32, error) {
	if m.openChannel != nil {
		return m.openChannel(ctx, peer, localSat, pushSat, private)
	}
	return chainhash.Hash{}, 0, nil
}
