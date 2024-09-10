package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func printQueryParams(val url.Values) {

	if len(val) <= 0 {
		return
	}

	fmt.Print("\n")
	fmt.Println("Query params:")
	for key, values := range val {
		if len(values) > 1 {
			fmt.Println("key:", key)
			fmt.Print("values: ")
			for _, val := range values {
				fmt.Printf("%s ", val)
			}
			fmt.Print("\n")
		} else {
			fmt.Printf("key: %s | value: %s\n", key, values[0])
		}
	}
}

func printRequestBody(body []byte) {
	var payload map[string]interface{}

	err := json.Unmarshal(body, &payload)

	if err != nil {
		return
	}

	if len(payload) <= 0 {
		return
	}

	fmt.Print("\n")
	fmt.Println("Request body:")

	for key, val := range payload {
		fmt.Printf("%v: %v\n", key, val)
	}
}

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("error reading body")
		}

		fmt.Println("-------------------Incoming request-------------------")
		fmt.Println("Request from IP:", r.RemoteAddr)
		fmt.Println("Method:", r.Method)
		fmt.Println("Requested URL:", r.URL.Path)
		printQueryParams(r.URL.Query())
		printRequestBody(bodyBytes)
		fmt.Println("-------------------Incoming request-------------------")

		// restoring the body as once i have read it here it is consumed as it is a stream of bytes
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	}
}
