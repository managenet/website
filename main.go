package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	port string = os.Getenv("PORT")
	csp string = "onnect-src 'self', script-src 'self', font-src https://fonts.gstatic.com"
	headers map[string]string = map[string]string{
		"Content-Security-Policy": csp,
		"X-Content-Security-Policy": csp,
		"X-Webkit-CSP": csp,
		"X-Frame-Options": "deny",
		"X-XSS-Protection": "1; mode=block",
		"X-Content-Type-Options": "nosniff",
		"Strict-Transport-Security": "max-age=16070400; includeSubDomains",
	}
)

func stampHeaders(w http.ResponseWriter, r *http.Request) {
	for name, value := range headers {
		w.Header().Set(name, value)
	}
}

func middlewareHandler(middleware func (http.ResponseWriter, *http.Request), handler http.Handler) func (http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		middleware(w, r) // Invoke the middleware (a func
		handler.ServeHTTP(w, r) // Invoke the actual handler
	}
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(
		middlewareHandler(stampHeaders,
			http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":" + port, router)
}