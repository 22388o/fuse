package fuse

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/lightning"
)

type lightningService interface {
	// Balance
	WalletBalance(ctx context.Context) (btcutil.Amount, error)

	// Payment
	AddInvoice(ctx context.Context, value lnwire.MilliSatoshi, memo string) (lightning.Invoice, error)
	PayInvoice(ctx context.Context, invoice lightning.Invoice) (lightning.PaymentResult, error)

	// Channels
	OpenChannel(ctx context.Context, addr lightning.LightningAddress, localSats, pushStats btcutil.Amount, private bool) (chainhash.Hash, uint32, error)
}
