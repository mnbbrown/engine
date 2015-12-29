package engine

import (
	log "github.com/Sirupsen/logrus"
	"github.com/dchest/uniuri"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, true)
				log.Debugf("%s", string(buf[0:stackSize]))
				JSON(rw, &J{
					"status_code": http.StatusInternalServerError,
					"message":     "Oops. Something went wrong",
				}, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(rw, req)
	})
}

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
