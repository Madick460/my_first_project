package middlewares

import (
	"errors"
	"github.com/Rasikrr/my_project/api"
	"log"
	"net/http"
)

type sessionKeyType string

const (
	SessionKey sessionKeyType = "session"
	authHeader string         = "Authorization"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("entering auth middleware")
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			api.SendError(w, http.StatusUnauthorized, errors.New("api key is empty"))
			return
		}
		next(w, r)
	}
}
