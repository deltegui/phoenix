package locomotive

import (
	"locomotive/injector"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//HTTPMethod is a typed version of HTTP methods
type HTTPMethod string

// Strong typed version for HTTP verbs
const (
	Get    HTTPMethod = "Get"
	Post   HTTPMethod = "Post"
	Delete HTTPMethod = "Delete"
)

type Mapping struct {
	Method   HTTPMethod
	Handler  http.HandlerFunc
	Endpoint string
}

type Controller interface {
	GetMappings() []Mapping
}

var router *mux.Router = mux.NewRouter()

func Map(root string, controllerBuilder injector.Builder) {
	controller := injector.CallBuilder(controllerBuilder).(Controller)
	middlewares := createMiddlewareChainWith(logMiddleware)
	mappings := controller.GetMappings()
	subRouter := router.PathPrefix(root).Subrouter()
	for _, mapping := range mappings {
		if mapping.Endpoint == "404" {
			router.NotFoundHandler = mapping.Handler
		} else {
			subRouter.HandleFunc(mapping.Endpoint, middlewares(mapping.Handler)).Methods(string(mapping.Method))
		}
	}
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("Request [%s] %s from %s\n", req.Method, req.RequestURI, req.RemoteAddr)
		next.ServeHTTP(writer, req)
	}
}

func createMiddlewareChainWith(chain ...middleware) middleware {
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

func startListening(address string) {
	log.Println("Listening on address: ", address)
	log.Println("You are ready to GO!")
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalln("Error while listening: ", err)
	}
}

func showEndpoints(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Println("Error while fetching route: ", err)
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			methods = []string{"ROOT ENDPOINT"}
		}
		log.Printf("%s Registered endpoint in %s", methods, path)
		return nil
	})
}
