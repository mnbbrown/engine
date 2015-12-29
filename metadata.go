package engine

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/dchest/uniuri"
	"net/http"
	"runtime"
	"time"
)

type RequestMetadata struct {
	RequestID string
	Method    string
	Path      string
	Status    int
	IP        string
	Size      int
	Latency   time.Duration
	StartTime time.Time
}

func (r *RequestMetadata) Fields() log.Fields {
	return log.Fields{
		"request_id": r.RequestID,
	}
}

func (r *RequestMetadata) Logger() *log.Entry {
	return log.WithFields(r.Fields())
}

type key int

const metadataCtxKey key = 0

func MetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		start := time.Now().UTC()

		metadata := &RequestMetadata{
			RequestID: uniuri.New(),
			StartTime: start,
			Method:    req.Method,
			Path:      req.URL.Path,
		}
		req.Header.Set("X-Request-Id", metadata.RequestID)
		rw.Header().Set("Request-Id", metadata.RequestID)

		metadata.IP = req.Header.Get("X-Real-IP")
		if metadata.IP == "" {
			metadata.IP = req.Header.Get("X-Forwarded-For")
			if metadata.IP == "" {
				metadata.IP = req.RemoteAddr
			}
		}

		// Set context
		GetContext(req).Set(metadataCtxKey, metadata)

		// Use StatusStoringRequest to keep a track of the request status.
		rw = NewResponseWriter(rw)

		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, true)
				metadata.Logger().Debugf("%s", string(buf[0:stackSize]))
				JSONError(rw, errors.New("Ooops. Something went wrong on our end"), http.StatusInternalServerError)
			}

			resp := rw.(*ResponseWriter)
			metadata.Status = resp.Status()
			metadata.Size = resp.Length()
			statusColor := colorForStatus(metadata.Status)
			metadata.Latency = time.Since(start)

			log.WithFields(metadata.Fields()).WithFields(log.Fields{
				"remote_ip": metadata.IP,
				"method":    metadata.Method,
				"path":      metadata.Path,
				"size":      metadata.Size,
				"latency":   metadata.Latency,
				"status":    metadata.Status,
			}).Printf("%s\t| %s | %s%d%s | %v | %s", metadata.Path, metadata.Method, statusColor, metadata.Status, reset, metadata.Latency, humanSize(metadata.Size))
		}()

		// Serve
		next.ServeHTTP(rw, req)
	})
}

func FromContext(ctx *Context) (*RequestMetadata, bool) {
	md, ok := ctx.Value(metadataCtxKey).(*RequestMetadata)
	return md, ok
}
