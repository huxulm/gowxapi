package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackdon/gowxapi/api"
	"github.com/jackdon/gowxapi/api/codesandbox"
	"github.com/jackdon/gowxapi/config"
	_ "github.com/jackdon/gowxapi/helper"
	_ "github.com/jackdon/gowxapi/seed"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/tools/playground/socket"

	// Imports so that go build/install automatically installs them.
	_ "golang.org/x/tour/pic"
	_ "golang.org/x/tour/tree"
	_ "golang.org/x/tour/wc"
)

const (
	basePkg    = "golang.org/x/tour"
	socketPath = "/socket"
)

var (
	httpListen = config.C.HttpListen
)

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
	if strings.HasSuffix(r.RequestURI, ".html") {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	}
	m.next.ServeHTTP(w, r)
	// We can
}

func main() {
	router := httprouter.New()
	router.GET("/login", api.Login)
	router.GET("/example/page", api.PageGoExample)
	router.POST("/example/detail/:id", api.GoExampleDetail)
	router.POST("/example/code/hl", api.HighLight)
	router.POST("/sandbox/list", codesandbox.ListSandBox)
	router.POST("/sandbox/lesson/page", codesandbox.LessonPaging)
	router.POST("/sandbox/section/detail/:id", codesandbox.GetLessonSectionDetail)
	router.POST("/sandbox/section/page/:lesson", codesandbox.LessonSectionPaging)

	host, port, err := net.SplitHostPort(httpListen)
	if err != nil {
		log.Fatal(err)
	}

	wsOrigin := strings.Split(config.C.WsOrigin, ",")
	origin := &url.URL{Scheme: wsOrigin[0], Host: wsOrigin[1]}
	router.Handler("GET", socketPath, socket.NewHandler(origin))
	router.ServeFiles("/static/*filepath", http.Dir("./"))
	httpAddr := host + ":" + port
	m := NewMiddleware(router, "middleware")

	log.Fatal(http.ListenAndServe(httpAddr, m))
}
