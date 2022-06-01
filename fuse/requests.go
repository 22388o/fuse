package fuse

import (
	"errors"
	"net/http"

	"github.com/mdedys/fuse/lightning"
)

type PayRequest struct {
	Request string `json:"id"`
}

func (a *PayRequest) Bind(r *http.Request) error {
	if a.Request == "" {
		return errors.New("missing required Request field")
	}
	return nil
}

type CreateInvoiceRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}

func (i *CreateInvoiceRequest) Bind(r *http.Request) error {
	if i.Amount <= 0 {
		return errors.New("invalid amount set for invoice")
	}
	return nil
}

type OpenChannelRequest struct {
	Addr        lightning.LightningAddress `json:"addr"`
	LocalAmount int64                      `json:"local_amount"`
	PushAmount  int64                      `json:"push_amount"`
}

func (ocr *OpenChannelRequest) Bind(r *http.Request) error {
	if ocr.Addr == "" {
		return errors.New("addr is required to open a channel")
	}
	if ocr.LocalAmount == 0 {
		return errors.New("must commit at least 1 sat to channel")
	}
	if ocr.PushAmount > ocr.LocalAmount {
		return errors.New("cannot push more sats then commited to channel")
	}
	return nil
}

type CreateLNURLPCodeRequest struct {
	MinSendable int64 `json:"min_sendable"`
	MaxSendable int64 `json:"max_sendable"`
}

func (l *CreateLNURLPCodeRequest) Bind(r *http.Request) error {
	if l.MaxSendable < l.MinSendable {
		return errors.New("max_sendable must be larger thenmin_sendable")
	}
	return nil
}
