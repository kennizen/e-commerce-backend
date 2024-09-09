package middlewares

import (
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-------------------Incoming request-------------------")
		fmt.Println("Request from IP", r.RemoteAddr)
		fmt.Println("Method", r.Method)
		fmt.Println("Requested URL", r.RequestURI)

		fmt.Println("-------------------Incoming request-------------------")

		next.ServeHTTP(w, r)
	}
}
