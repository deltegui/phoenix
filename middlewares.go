package phoenix

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("Request [%s] %s from %s\n", req.Method, req.RequestURI, req.RemoteAddr)
		next.ServeHTTP(writer, req)
	}
}

type csrfHandler struct {
	handler http.HandlerFunc
}

func (csrf csrfHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	csrf.handler(writer, req)
}

func NewCSRFMiddleware() Middleware {
	CSRF := csrf.Protect([]byte("adfaf"))
	return func(next http.HandlerFunc) http.HandlerFunc {
		return CSRF(csrfHandler{next}).ServeHTTP
	}
}

func createMiddlewareChainWith(chain []Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			last := final
			for i := len(chain) - 1; i >= 0; i-- {
				last = chain[i](last)
			}
			last(writer, req)
		}
	}
}
