package lnd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/lightning"
)

var (
	NullPreimage          = lntypes.Preimage{}
	ErrInvoiceAlreadyPaid = errors.New("invoice has already been paid")
)

type lnd interface {
	AddInvoice(ctx context.Context, in *invoicesrpc.AddInvoiceData) (lntypes.Hash, string, error)
	PayInvoice(ctx context.Context, invoice string, maxFee btcutil.Amount, outgoingChannel *uint64) chan lndclient.PaymentResult
	WalletBalance(ctx context.Context) (*lndclient.WalletBalance, error)
}

type LndClient struct {
	lnd     lnd
	network lightning.Network
	maxFee  btcutil.Amount
}

func (c LndClient) WalletBalance(ctx context.Context) (btcutil.Amount, error) {
	balance, err := c.lnd.WalletBalance(ctx)
	if err != nil {
		return btcutil.Amount(0), err
	}
	return balance.Confirmed, nil
}

func (c LndClient) Pay(ctx context.Context, invoice lightning.Invoice) (lightning.PaymentResult, error) {
	result := <-c.lnd.PayInvoice(ctx, invoice.Encoded, c.maxFee, nil)
	if result.Err != nil {
		return lightning.PaymentResult{}, result.Err
	}

	if result.Preimage == NullPreimage {
		return lightning.PaymentResult{}, ErrInvoiceAlreadyPaid
	}

	return lightning.PaymentResult{PreImage: result.Preimage, PaidFee: result.PaidFee}, nil
}

func (c LndClient) AddInvoice(ctx context.Context, value lnwire.MilliSatoshi, memo string) (lightning.Invoice, error) {

	data := &invoicesrpc.AddInvoiceData{Memo: memo, Value: value}

	_, encoded, err := c.lnd.AddInvoice(ctx, data)
	if err != nil {
		return lightning.Invoice{}, err
	}

	invoice, err := lightning.DecodeInvoice(encoded, c.network)
	if err != nil {
		return lightning.Invoice{}, err
	}

	return invoice, nil
}

func connect(address, network, macPath, tlsPath string) (lndclient.LightningClient, error) {

	cfg := &lndclient.LndServicesConfig{
		LndAddress:         address,
		Network:            lndclient.Network(network),
		CustomMacaroonPath: macPath,
		TLSPath:            tlsPath,
	}

	var lnd lndclient.LightningClient
	err := retry.Do(
		func() error {
			services, err := lndclient.NewLndServices(cfg)
			if err != nil {
				fmt.Printf("Failed to connect to LND: %s", err.Error())
				return err
			}
			lnd = services.Client
			return nil
		},
		retry.Attempts(10),
		retry.Delay(time.Duration(1)*time.Second),
	)

	return lnd, err
}

func NewClient(address, network, macPath, tlsPath string, maxFee btcutil.Amount) (*LndClient, error) {
	lnd, err := connect(address, network, macPath, tlsPath)
	if err != nil {
		return nil, err
	}
	return &LndClient{lnd: lnd, network: lightning.Network(network), maxFee: maxFee}, err
}
