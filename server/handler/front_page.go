package handler

import (
	"fmt"
	"net/http"
)

// Front is used to serve front page markup.
// 		GET /
// 		Responds: 200
func Front(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/api/v1/login">Google Log In</a></body></html>`
	fmt.Fprintf(w, html)
}
