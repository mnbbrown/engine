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
	enc.Encode(v)
}
