package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
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

func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Server error",
		Message:        err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func main() {
	// Connection string for a remote PostgreSQL database
	connStr := "postgres://api:awesomepassword@db:5432/api?sslmode=disable"

	// Open a connection to the remote database
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to the remote database!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Mount("/v1", AppRouter(db))
	http.ListenAndServe(":8000", r)
}

func AppRouter(db *sqlx.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gas, ok gas"))
	})
	r.Mount("/user", UserRouter(db))
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
	*User
	AccessToken string `json:"accessToken"`
}

type UserRegisterResponse struct {
	BaseResponse
	Data UserWithToken `json:"data"`
}

func UserRouter(db *sqlx.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok gas, ok gas"))
	})
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		UserRegistrationHandler(w, r, db)
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		UserLoginHandler(w, r, db)
	})
	return r
}

func isValidPassword(password string) bool {
	return len(password) >= 5 && len(password) <= 15
}

func isValidUserName(username string) bool {
	return len(username) >= 5 && len(username) <= 15
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

func generateToken(user *User) (string, error) {
	_, token, err := jwtauth.New("HS256", []byte("secret"), nil).Encode(map[string]interface{}{"user_id": user.ID})
	return token, err
}

func hashPassword(password string) (string, error) {
	saltLength, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))

	if err != nil {
		return "", err
	}

	if saltLength <= 0 {
		saltLength = 8
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), saltLength)
	return string(bytes), err
}

func comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func getUserByUserName(db *sqlx.DB, username string) (*User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	data := &User{}
	err := db.Get(data, query, username)
	return data, err
}

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	payload := &UserRegisterRequest{}

	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if !isValidPassword(payload.Password) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Password must be between 5 and 15 characters")))
		return
	}

	hashedPassword, err := hashPassword(payload.Password)

	if err != nil {
		render.Render(w, r, ErrServer(fmt.Errorf("Something wrong, please try again.")))
		return
	}

	token, err := generateToken(&User{ID: "1"})

	if err != nil {
		render.Render(w, r, ErrServer(fmt.Errorf("Error generating token, please try again.")))
		return
	}

	data := &UserWithToken{
		User:        &payload.User,
		AccessToken: token,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserRegisterResponse(data))
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	payload := &UserLoginRequest{}

	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if !isValidPassword(payload.Password) || !isValidUserName(payload.UserName) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid username or password")))
		return
	}

	user, err := getUserByUserName(db, payload.UserName)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid username or password")))
		return
	}

	if !comparePassword(user.Password, payload.Password) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid username or password")))
		return
	}

	token, err := generateToken(&User{ID: "1"})

	if err != nil {
		render.Render(w, r, ErrServer(fmt.Errorf("Error generating token, please try again.")))
		return
	}

	data := &UserWithToken{
		User:        user,
		AccessToken: token,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, ValidUserLoginResponse(data))
}
