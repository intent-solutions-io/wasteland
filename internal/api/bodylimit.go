package api

import "net/http"

// MaxBytesBody returns middleware that limits request body size to n bytes.
// Requests exceeding the limit receive a 413 Request Entity Too Large response.
func MaxBytesBody(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}
