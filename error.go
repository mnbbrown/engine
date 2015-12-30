package engine

import (
	"errors"
	"net/http"
)

func Error500(rw http.ResponseWriter) {
	JSONError(rw, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
}
