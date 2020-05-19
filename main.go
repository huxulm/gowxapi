package main

import (
	"log"
	"net/http"

	"github.com/jackdon/gowxapi/api"
	_ "github.com/jackdon/gowxapi/config"
	_ "github.com/jackdon/gowxapi/helper"
	_ "github.com/jackdon/gowxapi/seed"
	"github.com/julienschmidt/httprouter"
)

func init() {
}

// Middleware The type of our middleware consists of the original handler we want to wrap and a message
type Middleware struct {
	next    http.Handler
	message string
}

// NewMiddleware Make a constructor for our middleware type since its fields are not exported (in lowercase)
func NewMiddleware(next http.Handler, message string) *Middleware {
	return &Middleware{next: next, message: message}
}

// Our middleware handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We can modify the request here; for simplicity, we will just log a message
	log.Printf("Method: %s, URI: %s\n", r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	m.next.ServeHTTP(w, r)
	// We can
}

func main() {
	router := httprouter.New()
	router.GET("/login", api.Login)
	router.GET("/example/page", api.PageGoExample)
	router.POST("/example/detail/:id", api.GoExampleDetail)
	router.POST("/example/code/hl", api.HighLight)
	m := NewMiddleware(router, "")
	log.Fatal(http.ListenAndServe(":8080", m))
}
