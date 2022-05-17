package fuse

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mdedys/fuse/api"
	"github.com/mdedys/fuse/lnd"
)

type Fuse struct {
	lightning lightning
}

func (f Fuse) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	balance, err := f.lightning.WalletBalance(ctx)
	if err != nil {
		api.RespondWithError(w, 500, err)
	}

	api.RespondWithJSON(w, 200, balance)
}

func New(lightning lnd.LndClient) *chi.Mux {
	f := Fuse{
		lightning: lightning,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/balance", f.GetBalance)

	return r
}
