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

type Product struct {
	ID             string `json:"id,omitempty" db:"id"`
	Name           string `json:"name" db:"name"`
	Price          int    `json:"price" db:"price"`
	Condition      string `json:"condition" db:"condition"`
	IsPurchaseAble bool   `json:"isPurchaseAble" db:"is_purchase_able"`
	CreatedAt      string `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      string `json:"updatedAt,omitempty" db:"updated_at"`
}

type ProductCreateRequest struct {
}

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

func CreateProductHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	// payload := &ProductCreateRequest{}
	// if err := render.Decode(r, payload); err != nil {
	// 	render.Render(w, r, ErrInvalidRequest(err))
	// 	return
	// }
	// product := &Product{
	// 	ID:          ulid.Make().String(),
	// 	Name:        payload.Name,
	// 	Description: payload.Description,
	// 	Price:       payload.Price,
	// }
	// err := createProduct(db, product)
	// if err != nil {
	// 	render.Render(w, r, ErrServer(ParseDBErrorMessage(err)))
	// 	return
	// }
	// render.Status(r, http.StatusCreated)
	// render.Render(w, r, NewProductResponse(product))
}
