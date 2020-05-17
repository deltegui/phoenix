package locomotive

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deltegui/locomotive/injector"

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

type Controller interface {
	GetMappings() []CMapping
}

// Mapping represent a HTTP mapping for a Builder.
type Mapping struct {
	Method   HTTPMethod
	Endpoint string
	Builder  injector.Builder
}

// CMapping (Controller Mapping) represent a HTTP mapping for a http.HandlerFunc. Used mainly in controllers.
type CMapping struct {
	Method   HTTPMethod
	Endpoint string
	Handler  http.HandlerFunc
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
	mapper.cmap(CMapping{mapping.Method, mapping.Endpoint, controller})
}

func (mapper Mapper) MapRoot(controllerBuilder injector.Builder) {
	mapper.Map(Mapping{Get, "", controllerBuilder})
}

func (mapper Mapper) MapController(endpoint string, controllerBuilder injector.Builder) {
	controller := injector.CallBuilder(controllerBuilder).(Controller)
	mappings := controller.GetMappings()
	submapper := mapper.subMapperFrom(endpoint)
	for _, mapping := range mappings {
		submapper.cmap(mapping)
	}
}

func (mapper Mapper) MapRootController(controllerBuilder injector.Builder) {
	mapper.MapController("/", controllerBuilder)
}

func (mapper Mapper) subMapperFrom(endpoint string) Mapper {
	return Mapper{mapper.router.PathPrefix(endpoint).Subrouter()}
}

func (mapper Mapper) cmap(mapping CMapping) {
	middlewares := createMiddlewareChainWith(logMiddleware)
	mapper.router.HandleFunc(mapping.Endpoint, middlewares(mapping.Handler)).Methods(string(mapping.Method))
	mapper.router.HandleFunc(fmt.Sprintf("%s/", mapping.Endpoint), middlewares(mapping.Handler)).Methods(string(mapping.Method))
}

var mainMapper Mapper = Mapper{mux.NewRouter()}

func Map(mapping Mapping) {
	mainMapper.Map(mapping)
}

func MapRoot(controllerBuilder injector.Builder) {
	mainMapper.MapRoot(controllerBuilder)
}

func MapGroup(root string, createGroup func(mapper Mapper)) {
	createGroup(mainMapper.subMapperFrom(root))
}

func MapController(endpoint string, controllerBuilder injector.Builder) {
	mainMapper.MapController(endpoint, controllerBuilder)
}

func MapRootController(controllerBuilder injector.Builder) {
	mainMapper.MapRootController(controllerBuilder)
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
	if err := http.ListenAndServe(address, mainMapper.router); err != nil {
		log.Fatalln("Error while listening: ", err)
	}
}

func showEndpoints(router *mux.Router) {
	stringEndsWith := func(target, end string) bool {
		return strings.LastIndex(target, end) == len(target)-1
	}

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Println("Error while fetching route: ", err)
			return nil
		}
		if stringEndsWith(path, "/") && len(path) != 1 {
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
