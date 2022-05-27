package lnd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/routing/route"
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

	Connect(ctx context.Context, peer route.Vertex, host string, permanent bool) error
	ListPeers(ctx context.Context) ([]lndclient.Peer, error)
	OpenChannel(ctx context.Context, peer route.Vertex, localSat, pushSat btcutil.Amount, private bool) (*wire.OutPoint, error)
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

func (c LndClient) PayInvoice(ctx context.Context, invoice lightning.Invoice) (lightning.PaymentResult, error) {
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

func (c LndClient) ListPeers(ctx context.Context) ([]lightning.Peer, error) {
	_, err := c.lnd.ListPeers(ctx)
	if err != nil {
		return []lightning.Peer{}, err
	}
	return []lightning.Peer{}, nil
}

func (c LndClient) ConnectPeer(ctx context.Context, peer lightning.Vertex, host string) error {
	return errors.New("Not Implemented")
}

func (c LndClient) OpenChannel(ctx context.Context, peer lightning.Vertex, localSat, pushSat btcutil.Amount, private bool) (chainhash.Hash, uint32, error) {
	return chainhash.Hash{}, 0, errors.New("Not Implemented")
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
