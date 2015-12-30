package metadata

import (
	"time"

	log "github.com/Sirupsen/logrus"
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
