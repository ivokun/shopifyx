package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID        string     `json:"-" db:"id"`
	Name      string     `json:"name" sql:"name"`
	UserName  string     `json:"username" sql:"username"`
	Password  string     `json:"-" sql:"password"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
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

// User utils
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
	_, token, err := TokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
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

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// User database operations
func getUserByUserName(db *sqlx.DB, username string) (*User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	data := &User{}
	err := db.Get(data, query, username)
	return data, err
}

func createUser(db *sqlx.DB, user *User) error {
	_, err := db.NamedExec(`INSERT INTO users (id, name, username, password) VALUES (:id, :name, :username, :password)`, user)
	return err
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
		render.Render(w, r, ErrServer(fmt.Errorf("Something wrong, please try again."), http.StatusInternalServerError))
		return
	}

	user := &User{
		ID:       ulid.Make().String(),
		Name:     payload.User.Name,
		UserName: payload.User.UserName,
		Password: hashedPassword,
	}

	err = createUser(db, user)

	if err != nil {
		render.Render(w, r, ErrServer(ParseDBErrorMessage(err)))
		return
	}

	token, err := generateToken(user)

	if err != nil {
		render.Render(w, r, ErrServer(fmt.Errorf("Error generating token, please try again."), http.StatusInternalServerError))
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
		render.Render(w, r, ErrServer(ParseDBErrorMessage(err)))
		return
	}

	if !comparePassword(user.Password, payload.Password) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid username or password")))
		return
	}

	token, err := generateToken(user)

	if err != nil {
		render.Render(w, r, ErrServer(fmt.Errorf("Error generating token, please try again."), http.StatusInternalServerError))
		return
	}

	data := &UserWithToken{
		User:        user,
		AccessToken: token,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, ValidUserLoginResponse(data))
}
