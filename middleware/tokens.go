package middleware

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/seanburman/game-ws-server/config"
	"github.com/seanburman/game-ws-server/services"
)

func MiddlewareVerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := services.AuthService{}.StripToken(r)
		if err != nil {
			log.Println("cannot parse token from header")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("invalid access token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("unauthorized"))
			}
			return []byte(config.Env().SECRET), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			// TODO: implement proper error package
			if _, err = w.Write([]byte("invalid access token")); err != nil {
				log.Println("error writing to client:", err)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
