package lightning

import (
	"context"

	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
)

type PaymentResult struct {
	PreImage [32]byte
	PaidFee  btcutil.Amount
}

// AddInvoice creates a bolt11 invoice
func (l LightningClient) AddInvoice(ctx context.Context, value lnwire.MilliSatoshi, memo string, hhash []byte) (Invoice, error) {
	return l.provider.AddInvoice(ctx, value, memo, hhash)
}

// PayInvoice pays lightning invoice
func (l LightningClient) PayInvoice(ctx context.Context, invoice Invoice) (PaymentResult, error) {
	return l.provider.PayInvoice(ctx, invoice)
}

// WalletBalance retrieves the wallet balance on lightning node
func (l LightningClient) WalletBalance(ctx context.Context) (btcutil.Amount, error) {
	return l.provider.WalletBalance(ctx)
}
