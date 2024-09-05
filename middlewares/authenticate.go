package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kennizen/e-commerce-backend/lib"
	"github.com/kennizen/e-commerce-backend/utils"
)

type ContextKey string

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
			return
		}

		bearerToken := strings.Split(token, " ")

		if bearerToken[0] != "Bearer" && bearerToken[1] == "" {
			utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
			return
		}

		claims, isValid := lib.ValidateToken(bearerToken[1], os.Getenv("JWT_TOKEN_SECRET"))

		fmt.Println("claims", claims)

		if !isValid {
			utils.SendMsg("Invalid token", http.StatusUnauthorized, w)
		} else {
			ctx := context.WithValue(r.Context(), ContextKey("userID"), claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
