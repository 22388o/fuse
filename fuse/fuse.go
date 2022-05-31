package fuse

import (
	"net/http"

	"github.com/btcsuite/btcutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/api"
	"github.com/mdedys/fuse/lightning"
)

type Fuse struct {
	lightning lightningService
	network   lightning.Network
}

// repsondWithJSON creates a successful json response
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, v render.Renderer) {
	render.SetContentType(render.ContentTypeJSON)
	render.Status(r, code)
	render.Render(w, r, v)
}

// GetBalance retrieves the wallet balance on the node
func (f Fuse) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	balance, err := f.lightning.WalletBalance(ctx)
	if err != nil {
		api.RespondWithError(w, 500, err)
	}

	api.RespondWithJSON(w, 200, balance)
}

// Pay pays a lightning invoice
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

	result, err := f.lightning.PayInvoice(ctx, invoice)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewPayResponse(result))
}

// CreateInvoice creates a bolt11 invoice
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

	respondWithJSON(w, r, http.StatusCreated, NewInvoiceResponse(invoice))
}

// OpenChannel opens a lightning network channel
func (f Fuse) OpenChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := &OpenChannelRequest{}
	if err := render.Bind(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	hash, idx, err := f.lightning.OpenChannel(ctx, payload.Addr, btcutil.Amount(payload.LocalAmount), btcutil.Amount(payload.PushAmount), false)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	respondWithJSON(w, r, http.StatusCreated, NewOpenChannelResponse(hash, idx))
}

func New(lightning lightningService, network lightning.Network) *chi.Mux {
	f := Fuse{
		lightning: lightning,
		network:   network,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/balance", f.GetBalance)
	r.Post("/pay", f.Pay)
	r.Post("/invoices", f.CreateInvoice)
	r.Post("/channels", f.OpenChannel)

	return r
}
