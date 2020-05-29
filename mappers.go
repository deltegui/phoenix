package phoenix

import (
	"fmt"
	"net/http"

	"github.com/deltegui/phoenix/injector"

	"github.com/gorilla/mux"
)

// HTTPMethod is a typed version of HTTP methods
type HTTPMethod string

// Strong typed version for HTTP verbs
const (
	Get     HTTPMethod = "Get"
	Post    HTTPMethod = "Post"
	Delete  HTTPMethod = "Delete"
	Head    HTTPMethod = "Head"
	Put     HTTPMethod = "Put"
	Connect HTTPMethod = "Connect"
	Options HTTPMethod = "Options"
	Trace   HTTPMethod = "Trace"
	Patch   HTTPMethod = "Delete"
)

// Mapping represent a HTTP mapping for a Builder.
type Mapping struct {
	Method      HTTPMethod
	Endpoint    string
	Builder     injector.Builder
	middlewares []Middleware
}

func (mapping *Mapping) Use(middleware Middleware) *Mapping {
	mapping.middlewares = append(mapping.middlewares, middleware)
	return mapping
}

type Mapper struct {
	router *mux.Router
}

func (mapper Mapper) Map(mapping Mapping) {
	controller := injector.CallBuilder(mapping.Builder).(http.HandlerFunc)
	if mapping.Endpoint == "404" {
		mapper.router.NotFoundHandler = controller
		return
	}
	mapsWithLog := append(mapping.middlewares, logMiddleware)
	middlewares := createMiddlewareChainWith(mapsWithLog)
	mapper.router.HandleFunc(mapping.Endpoint, middlewares(controller)).Methods(string(mapping.Method))
	mapper.router.HandleFunc(fmt.Sprintf("%s/", mapping.Endpoint), middlewares(controller)).Methods(string(mapping.Method))
}

func (mapper Mapper) MapAll(mappings []Mapping) {
	for _, mapping := range mappings {
		mapper.Map(mapping)
	}
}

func (mapper Mapper) MapRoot(controllerBuilder injector.Builder) {
	mapper.Map(Mapping{
		Method:   Get,
		Endpoint: "",
		Builder:  controllerBuilder,
	})
}

func (mapper Mapper) MapGroup(root string, createGroup func(mapper Mapper)) {
	createGroup(mapper.subMapperFrom(root))
}

func (mapper Mapper) subMapperFrom(endpoint string) Mapper {
	return Mapper{
		router: mapper.router.PathPrefix(endpoint).Subrouter(),
	}
}
