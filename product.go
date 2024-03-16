package main

import (
	// "fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	// "github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"net/http"
	// "os"
	// "strconv"
	// "time"
)

func ProductRouter(db *sqlx.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gas, ok gas"))
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(TokenAuth))
		r.Use(jwtauth.Authenticator(TokenAuth))
		r.Get("/here", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected"))
		})
	})

	return r
}
