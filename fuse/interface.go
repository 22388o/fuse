package fuse

import (
	"context"

	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/lightning"
)

type lightningService interface {
	WalletBalance(ctx context.Context) (btcutil.Amount, error)
	AddInvoice(ctx context.Context, value lnwire.MilliSatoshi, memo string) (lightning.Invoice, error)
	Pay(ctx context.Context, invoice lightning.Invoice) (lightning.PaymentResult, error)
}
