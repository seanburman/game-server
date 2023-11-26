package services

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/seanburman/game-ws-server/config"
	"github.com/seanburman/game-ws-server/server"
)

type (
	AuthService struct {
		*server.ServerContext
	}
	JWTClaims struct {
		jwt.RegisteredClaims
		UserID string `json:"user_id,omitempty"`
	}
)

func NewAuthService(ctx *server.ServerContext) *AuthService {
	return &AuthService{ServerContext: ctx}
}

func (as AuthService) NewToken(userId string) string {
	// Set custom claims
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
		},
		UserID: userId,
	}
	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Encode token
	t, err := token.SignedString([]byte(config.Env().SECRET))
	if err != nil {
		panic(err)
	}
	return "Bearer " + t
}

func (as AuthService) StripToken(r *http.Request) (string, error) {
	authorization := r.Header["Authorization"]
	if len(authorization) == 0 {
		return "", errors.New("authorization header missing")
	}

	bearerToken := strings.Split(authorization[0], " ")
	return bearerToken[1], nil
}
