package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtS string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtS = cookie.Value
			} else {
				http.Error(w, "cookie is empty", http.StatusUnauthorized)
				return
			}
			var valid bool
			token, err := jwt.Parse(jwtS, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			valid = token.Valid
			if err != nil || !valid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
			hashFromToken, ok := claims["hash"].(string)
			if !ok {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
			hash := sha256.Sum256([]byte(pass))
			hashString := hex.EncodeToString(hash[:])
			if hashFromToken != hashString {
				http.Error(w, "cookie is not correct", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
