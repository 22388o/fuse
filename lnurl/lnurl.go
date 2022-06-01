package lnurl

import (
	"github.com/btcsuite/btcutil/bech32"
	"github.com/fiatjaf/go-lnurl"
)

func CreateBech32Code(url string) (string, error) {
	bytes := []byte(url)

	bits, err := bech32.ConvertBits(bytes, 8, 5, true)
	if err != nil {
		return "", err
	}

	code, err := lnurl.Encode("lnurl", bits)
	if err != nil {
		return "", err
	}
	return code, err
}
