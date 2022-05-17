package lnd

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/lndclient"
)

type lnd interface {
	WalletBalance(ctx context.Context) (*lndclient.WalletBalance, error)
}

type LndClient struct {
	lnd lnd
}

func (c LndClient) WalletBalance(ctx context.Context) (btcutil.Amount, error) {
	balance, err := c.lnd.WalletBalance(ctx)
	if err != nil {
		return btcutil.Amount(0), err
	}
	return balance.Confirmed, nil
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

func NewClient(address, network, macPath, tlsPath string) (*LndClient, error) {
	lnd, err := connect(address, network, macPath, tlsPath)
	if err != nil {
		return nil, err
	}
	return &LndClient{lnd}, err
}
