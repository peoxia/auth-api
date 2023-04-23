// Package handler contains HTTP handlers.
package handler

import "net/http"

// Healthz is used for application health monitoring.
// 		GET /_healthz
// 		Responds: 200
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
