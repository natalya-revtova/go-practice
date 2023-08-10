package hello

import (
	"net/http"

	"github.com/go-chi/render"
)

func SayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "hello")
	}
}
