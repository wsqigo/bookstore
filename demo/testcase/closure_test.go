package testcase

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	http.HandleFunc("/hello", timed(hello))
	http.ListenAndServe(":3000", nil)
}

func timed(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		f(w, r)
		fmt.Println("The request took", time.Since(start))
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello!</h1>")
}
