package phoenix

import (
	"log"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("Request [%s] %s from %s\n", req.Method, req.RequestURI, req.RemoteAddr)
		next.ServeHTTP(writer, req)
	}
}

func createMiddlewareChainWith(chain []Middleware) Middleware {
	chain = append(chain, logMiddleware)
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
