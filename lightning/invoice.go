package lightning

import (
	"errors"

	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/zpay32"
	"github.com/rs/zerolog/log"
)

type Invoice struct {
	Encoded string
	Decoded zpay32.Invoice
}

func DecodeInvoice(encoded string, network Network) (Invoice, error) {

	log.Info().Msg(string(network))

	cp, err := network.ChainParams()
	if err != nil {
		return Invoice{}, errors.New("unable to decode network")
	}

	decoded, err := zpay32.Decode(encoded, cp)
	if err != nil {
		return Invoice{}, errors.New("unable to decode bolt11 invoice")
	}

	invoice := Invoice{
		Encoded: encoded,
		Decoded: *decoded,
	}

	return invoice, nil
}

type PaymentResult struct {
	PreImage [32]byte
	PaidFee  btcutil.Amount
}
