package lnurl

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/fiatjaf/go-lnurl"
)

const TagPayRequst = "payRequest"

type MetadataImage struct {
	DataURI string
	Bytes   []byte
}

type Metadata struct {
	Description string
	Image       MetadataImage
}

func GetImage() (MetadataImage, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return MetadataImage{}, err
	}

	f, err := os.Open(filepath.Join(cwd, "../../lnurl", "ln_bolt.png"))
	if err != nil {
		return MetadataImage{}, err
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return MetadataImage{}, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return MetadataImage{}, err
	}

	bytes := buf.Bytes()
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return MetadataImage{DataURI: encoded, Bytes: bytes}, nil
}

func (c Metadata) Encode() string {
	description := [2]string{"text/plain", c.Description}
	image := [2]string{"image/png:base64", c.Image.DataURI}
	metadata := [2][2]string{description, image}
	encoded, _ := json.Marshal(metadata)
	return string(encoded)
}

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
