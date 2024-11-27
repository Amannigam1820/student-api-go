package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Amannigam1820/student-api-go/internal/utils/response"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("student_api_go") // Replace with a secure key

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(r.Cookies())

		c, err := r.Cookie("token")
		//fmt.Println(c, err)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("bad request")))
			return
		}
		tokenStr := c.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
