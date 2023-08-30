package internalhttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func parseBody(r *http.Request, body any) error {
	err := render.DecodeJSON(r.Body, body)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("request body is empty")
	}
	if err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}
	return nil
}

func parseID(r *http.Request) string {
	return chi.URLParam(r, "id")
}
