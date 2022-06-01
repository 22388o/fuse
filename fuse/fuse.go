package fuse

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/mdedys/fuse/api"
	"github.com/mdedys/fuse/lightning"
	"github.com/mdedys/fuse/lnurl"
)

type Fuse struct {
	lightning lightningService
	network   lightning.Network
	store     inMemoryStore
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

	invoice, err := f.lightning.AddInvoice(ctx, lnwire.NewMSatFromSatoshis(btcutil.Amount(payload.Amount)), payload.Memo, []byte{})
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

func (f Fuse) ListChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channels, err := f.lightning.ListChannels(ctx, false, false)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	respondWithJSON(w, r, http.StatusOK, NewListChannelsResponse(channels))
}

func (f Fuse) CreateLNURLPCode(w http.ResponseWriter, r *http.Request) {

	payload := &CreateLNURLPCodeRequest{}
	if err := render.Bind(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	url := fmt.Sprintf("http://localhost:1100/lnurlp?id=%v&min=%v&max=%v", lnurl.RandomK1(), payload.MinSendable, payload.MaxSendable)
	code, err := lnurl.CreateBech32Code(url)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Render(w, r, NewLNURLPCodeResponse(code))
}

func (f Fuse) HandleLNURLP(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		render.Render(w, r, ErrLNURLRequestError(errors.New("missing payment request identifier")))
		return
	}

	min, err := strconv.ParseInt(r.URL.Query().Get("min"), 10, 64)
	if err != nil {
		render.Render(w, r, ErrLNURLRequestError(err))
		return
	}

	max, err := strconv.ParseInt(r.URL.Query().Get("max"), 10, 64)
	if err != nil {
		render.Render(w, r, ErrLNURLRequestError(err))
		return
	}

	img, err := lnurl.GetImage()
	if err != nil {
		render.Render(w, r, ErrLNURLRequestError(err))
		return
	}

	metadata := lnurl.Metadata{
		Description: "Fuse LNURL Payment",
		Image:       img,
	}

	f.store.Set(id, metadata.Encode())

	render.Render(w, r, NewLNURLPResponse(fmt.Sprintf("http://localhost:1100/lnurlp/callback/%s", id), min, max, metadata, lnurl.TagPayRequst))
}

func (f Fuse) HandleLNURLPCallback(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	id := chi.URLParam(r, "requestID")

	amount, err := strconv.ParseInt(r.URL.Query().Get("amount"), 10, 64)
	if err != nil {
		render.Render(w, r, ErrLNURLRequestError(err))
		return
	}

	metadata := f.store.Get(id)
	hhash := sha256.Sum256([]byte(metadata))

	invoice, err := f.lightning.AddInvoice(ctx, lnwire.MilliSatoshi(amount), "", hhash[:])
	if err != nil {
		render.Render(w, r, ErrLNURLRequestError(err))
		return
	}

	render.Render(w, r, NewLNURLPCallbackResponse(invoice))
}

func New(lightning lightningService, network lightning.Network, store inMemoryStore) *chi.Mux {
	f := Fuse{
		lightning: lightning,
		network:   network,
		store:     store,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/balance", f.GetBalance)
	r.Post("/pay", f.Pay)
	r.Post("/invoices", f.CreateInvoice)

	r.Route("/channels", func(r chi.Router) {
		r.Post("/", f.OpenChannel)
		r.Get("/", f.ListChannels)
	})

	r.Route("/lnurlp", func(r chi.Router) {
		r.Post("/", f.CreateLNURLPCode)
		r.Get("/", f.HandleLNURLP)
		r.Get("/callback/{requestID}", f.HandleLNURLPCallback)
	})

	return r
}
