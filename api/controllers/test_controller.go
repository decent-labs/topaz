package controllers

import (
	"net/http"
)

// TestController is our "ping"
func TestController(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Write([]byte("Hello, World!"))
}
