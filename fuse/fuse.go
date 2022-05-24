package fuse

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/btcsuite/btcutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/api"
	"github.com/mdedys/fuse/lightning"
	"github.com/mdedys/fuse/lnd"
)

type Fuse struct {
	lightning lightningService
	network   lightning.Network
}

func (f Fuse) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	balance, err := f.lightning.WalletBalance(ctx)
	if err != nil {
		api.RespondWithError(w, 500, err)
	}

	api.RespondWithJSON(w, 200, balance)
}

type PayRequest struct {
	Request string `json:"id"`
}

func (a *PayRequest) Bind(r *http.Request) error {
	if a.Request == "" {
		return errors.New("missing required Request field")
	}
	return nil
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

func (f Fuse) Pay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := &PayRequest{}
	if err := render.Bind(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	invoice, err := lightning.DecodeInvoice(payload.Request, f.network)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	result, err := f.lightning.Pay(ctx, invoice)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewPayResponse(result))
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

type CreateInvoiceResponse struct {
	Invoice string `json:"invoice"`
}

func (cir *CreateInvoiceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewInvoiceResponse(i lightning.Invoice) *CreateInvoiceResponse {
	return &CreateInvoiceResponse{Invoice: i.Encoded}
}

func (f Fuse) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := &CreateInvoiceRequest{}
	if err := render.Bind(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	invoice, err := f.lightning.AddInvoice(ctx, lnwire.NewMSatFromSatoshis(btcutil.Amount(payload.Amount)), payload.Memo)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewInvoiceResponse(invoice))
}

func New(lightning lnd.LndClient, network lightning.Network) *chi.Mux {
	f := Fuse{
		lightning: lightning,
		network:   network,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/balance", f.GetBalance)
	r.Post("/pay", f.Pay)
	r.Post("/invoices", f.CreateInvoice)

	return r
}
