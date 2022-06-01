package fuse

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type ErrReponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrReponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrReponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) render.Renderer {
	log.Error().Err(err).Msg("Internal server error.")
	return &ErrReponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

type LNURLErrorResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (e *LNURLErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusBadRequest)
	return nil
}

func ErrLNURLRequestError(reason error) render.Renderer {
	return &LNURLErrorResponse{Status: "ERROR", Reason: reason.Error()}
}
