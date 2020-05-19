package api

import "net/http"

// RespError return http status code 204
func RespError(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusNoContent)
}
