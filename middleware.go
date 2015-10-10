package engine

import (
	"github.com/dchest/uniuri"
	"log"
	"net/http"
	"strings"
	"time"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		id := uniuri.NewLen(uniuri.UUIDLen)
		GetContext(req).Set("Request-Id", id)
		req.Header.Set("X-Request-Id", id)
		next.ServeHTTP(rw, req)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		log.Printf("%s %d", req.URL.Path, time.Now().Sub(start).Nanoseconds())
	})
}

func CORSAcceptAll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
			rw.WriteHeader(200)
			return
		} else {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Headers", "*")
			next.ServeHTTP(rw, req)
		}
	})
}

type CORSConfig struct {
	AllowedMethods []string
	AllowedOrigins []string
	AllowedHeaders []string
}

var AllowAllConfig = &CORSConfig{
	AllowedMethods: []string{"*"},
	AllowedOrigins: []string{"*"},
	AllowedHeaders: []string{"Authorization", "X-Requested-With"},
}

func CORSMiddleware(config *CORSConfig) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch {
			case req.Method == "OPTIONS":
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				rw.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ", "))
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				rw.WriteHeader(200)
				return
			default:
				rw.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				rw.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ", "))
				rw.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				next.ServeHTTP(rw, req)
			}
		})
	}
}
