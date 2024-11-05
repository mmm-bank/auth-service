package http

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mmm-bank/auth-service/models"
	"github.com/mmm-bank/auth-service/storage"
	"net/http"
	"time"
)

type AuthService struct {
	userRepo    storage.UserRepo
	sessionRepo storage.SessionRepo
	secretKey   []byte
}

func NewAuthService(userRepo storage.UserRepo, sessionRepo storage.SessionRepo, secretKey string) *AuthService {
	return &AuthService{
		userRepo,
		sessionRepo,
		[]byte(secretKey),
	}
}

func (a *AuthService) generateToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) registerHandler(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := a.userRepo.AddUser(user); err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}
}

func (a *AuthService) loginHandler(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := a.userRepo.AuthUser(user)
	if err != nil {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	token, err := a.generateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	session := models.Session{
		Token:     token,
		UserID:    user.ID.String(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = a.sessionRepo.AddSession(session)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func CreateAndRunServer(t *AuthService, addr string) error {
	r := chi.NewRouter()
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", t.registerHandler)
		r.Post("/login", t.loginHandler)
	})
	return http.ListenAndServe(addr, r)
}
