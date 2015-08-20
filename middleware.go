package engine

import (
	"github.com/dchest/uniuri"
	"net/http"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		id := uniuri.NewLen(uniuri.UUIDLen)
		GetContext(req).Set("Request-Id", id)
		req.Header.Set("X-Request-Id", id)
		next.ServeHTTP(rw, req)
	})
}

func CORSAcceptAll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.WriteHeader(200)
			return
		} else {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(rw, req)
		}
	})
}
