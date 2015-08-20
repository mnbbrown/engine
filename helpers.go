package engine

import (
	"encoding/json"
	"net/http"
)

type J map[string]interface{}

func JSON(rw http.ResponseWriter, v interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	enc := json.NewEncoder(rw)
	switch v.(type) {
	case error:
		enc.Encode(&J{"error": v.(error).Error(), "status_code": code})
	default:
		enc.Encode(v)
	}
}
