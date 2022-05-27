package lightning

import (
	"errors"

	"github.com/lightningnetwork/lnd/zpay32"
)

type Invoice struct {
	Encoded string
	Decoded zpay32.Invoice
}

// DecodeInvoice takes in a encoded bolt11 invoice and decodes the invoice into its parts
func DecodeInvoice(encoded string, network Network) (Invoice, error) {

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
