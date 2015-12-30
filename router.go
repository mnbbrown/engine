package engine

import (
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mnbbrown/logger"
	"golang.org/x/net/context"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Router struct {
	mux          *httprouter.Router
	absolutePath string
	middleware   []MiddlewareFunc
	logger       *logger.Logger
}

func (r *Router) ListMiddleware() (mi []string) {
	for _, m := range r.middleware {
		mi = append(mi, runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name())
	}
	return mi
}

func NewRouter() *Router {
	r := httprouter.New()
	return &Router{mux: r}
}

func (r *Router) UseLogger(l *logger.Logger) {
	r.logger = l
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(rw, req)
}

func (r *Router) SubRouter(relativePath string, middleware ...MiddlewareFunc) *Router {
	middleware = append(r.middleware, middleware...)
	sr := &Router{
		mux:          r.mux,
		absolutePath: relativePath,
		middleware:   middleware,
	}
	if r.logger != nil {
		sr.UseLogger(r.logger)
	}
	return sr
}

func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) Static(relativePath, root string) {
	absolutePath := r.calculateAbsolutePath(relativePath)
	absolutePath = path.Join(absolutePath, "/*filepath")

	r.mux.ServeFiles(absolutePath, http.Dir(root))
}

func (r *Router) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return r.absolutePath
	}

	absolutePath := path.Join(r.absolutePath, relativePath)
	if strings.HasSuffix(relativePath, "/") && !strings.HasSuffix(absolutePath, "/") {
		return absolutePath + "/"
	}
	return absolutePath

}

func (r *Router) UseHandler(handler http.Handler) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(rw, req)
			next.ServeHTTP(rw, req)
		})
	})
}

func (r *Router) Next(handler http.Handler) {
	r.mux.NotFound = http.HandlerFunc(handler.ServeHTTP)
}

func (r *Router) Handle(method, path string, handler http.Handler, middleware ...MiddlewareFunc) {
	absolutePath := r.calculateAbsolutePath(path)
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	r.mux.Handle(method, absolutePath, wrap(handler))
}

func (r *Router) HandleFunc(method, path string, handler func(http.ResponseWriter, *http.Request), middleware ...MiddlewareFunc) {
	r.Handle(method, path, http.HandlerFunc(handler), middleware...)
}

// Get registers a GET handler for the given path.
func (r *Router) Get(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("GET", path, handler, middleware...)
}

func (r *Router) Head(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("HEAD", path, handler, middleware...)
}

// Put registers a PUT handler for the given path.
func (r *Router) Put(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("PUT", path, handler, middleware...)
}

// Post registers a POST handler for the given path.
func (r *Router) Post(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("POST", path, handler, middleware...)
}

// Patch registers a PATCH handler for the given path.
func (r *Router) Patch(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("PATCH", path, handler, middleware...)
}

// Delete registers a DELETE handler for the given path.
func (r *Router) Delete(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("DELETE", path, handler, middleware...)
}

// Options registers a OPTIONS handler for the given path.
func (r *Router) Options(path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) {
	r.HandleFunc("OPTIONS", path, handler, middleware...)
}

type paramsContextWrapper struct {
	context.Context
	httprouter.Params
}

func wrap(handler http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := GetContext(req)
		ctx.Params = params
		handler.ServeHTTP(rw, req)
	}
}

type ResponseWriter struct {
	status int
	size   int
	http.ResponseWriter
}

func (s *ResponseWriter) Status() int {
	return s.status
}

func (s *ResponseWriter) Length() int {
	return s.size
}

func (s *ResponseWriter) Header() http.Header {
	return s.ResponseWriter.Header()
}

func (s *ResponseWriter) Write(data []byte) (int, error) {
	if s.status == 0 {
		s.WriteHeader(200)
	}
	n, err := s.ResponseWriter.Write(data)
	s.size += n
	return n, err
}

func (s *ResponseWriter) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: rw, status: 200}
}
