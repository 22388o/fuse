package fuse

import (
	"context"

	"github.com/btcsuite/btcutil"
)

type lightning interface {
	WalletBalance(ctx context.Context) (btcutil.Amount, error)
}
