package fuse

import (
	"encoding/hex"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/mdedys/fuse/lightning"
)

type CreateInvoiceResponse struct {
	Invoice string `json:"invoice"`
}

func (cir *CreateInvoiceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewInvoiceResponse(i lightning.Invoice) *CreateInvoiceResponse {
	return &CreateInvoiceResponse{Invoice: i.Encoded}
}

type PayResponse struct {
	PreImage string  `json:"preimage"`
	PaidFee  float64 `json:"paid_fee"`
}

func (pr *PayResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPayResponse(result lightning.PaymentResult) *PayResponse {
	return &PayResponse{
		PreImage: hex.EncodeToString(result.PreImage[:]),
		PaidFee:  result.PaidFee.ToUnit(btcutil.AmountSatoshi),
	}
}

type OpenChannelResponse struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

func (ocr *OpenChannelResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOpenChannelResponse(hash chainhash.Hash, index uint32) *OpenChannelResponse {
	return &OpenChannelResponse{
		Hash:  hash.String(),
		Index: index,
	}
}

type Channel struct {
	ID            int    `json:"id"`
	LocalBalance  int    `json:"local_balance"`
	RemoteBalance int    `json:"remote_balance"`
	Capacity      int    `json:"capacity"`
	Active        bool   `json:"active"`
	Private       bool   `json:"private"`
	RemotePubkey  string `json:"remote_pubkey"`
}

type ListChannelsResponse struct {
	Channels []Channel `json:"channels"`
}

func (lcr *ListChannelsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewListChannelsResponse(channels []lightning.Channel) *ListChannelsResponse {
	var chans []Channel

	for _, c := range channels {
		pubkey := hex.EncodeToString(c.RemotePubkey[:])

		chans = append(chans, Channel{
			ID:            int(c.ID),
			LocalBalance:  int(c.LocalBalance),
			RemoteBalance: int(c.RemoteBalance),
			Capacity:      int(c.Capacity),
			Active:        c.Active,
			Private:       c.Private,
			RemotePubkey:  pubkey,
		})
	}

	return &ListChannelsResponse{chans}
}
