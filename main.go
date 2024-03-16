package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"time"
	// "github.com/oklog/ulid/v2"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message    string `json:"message"` // user-level status message
	StatusText string `json:"status"`  // user-level status message
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		Message:        err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Mount("/v1", AppRouter())
	http.ListenAndServe(":8000", r)
}

func AppRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gas, ok gas"))
	})
	r.Mount("/user", UserRouter())
	return r
}

type BaseResponse struct {
	Message string `json:"message"`
}

// Auth section

type User struct {
	ID        string    `json:"-" db:"id"`
	Name      string    `json:"name" sql:"name"`
	UserName  string    `json:"username" sql:"username"`
	Password  string    `json:"-" sql:"password"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`
}

type UserRegisterRequest struct {
	User
	Password string `json:"password"`
}

type UserLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserWithToken struct {
	User
	AccessToken string `json:"accessToken"`
}

type UserRegisterResponse struct {
	BaseResponse
	Data UserWithToken `json:"data"`
}

func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gas, ok gas"))
	})
	r.Post("/register", userRegistrationHandler)
	r.Post("/login", userLoginHandler)
	return r
}

func isValidPassword(password string) bool {
	return len(password) >= 5 && len(password) <= 15
}

func isValidUserName(username string) bool {
	return len(username) >= 5 && len(username) <= 50
}

func NewUserRegisterResponse(data *UserWithToken) *UserRegisterResponse {
	resp := &UserRegisterResponse{
		BaseResponse: BaseResponse{
			Message: "User registered successfully",
		},
		Data: *data,
	}

	return resp
}

func ValidUserLoginResponse(data *UserWithToken) *UserRegisterResponse {
	resp := &UserRegisterResponse{
		BaseResponse: BaseResponse{
			Message: "User logged in successfully",
		},
		Data: *data,
	}

	return resp
}

func (rd *UserRegisterResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UserRegisterRequest{}

	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if !isValidPassword(payload.Password) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Password must be between 5 and 15 characters")))
		return
	}

	data := &UserWithToken{
		User:        payload.User,
		AccessToken: "token",
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserRegisterResponse(data))
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &UserLoginRequest{}

	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if !isValidPassword(payload.Password) && !isValidUserName(payload.UserName) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid username or password")))
		return
	}

	data := &UserWithToken{
		User: User{
			UserName: payload.UserName,
			Name:     "John Doe",
		},
		AccessToken: "token",
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, ValidUserLoginResponse(data))
}
