package phoenix

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/deltegui/phoenix/injector"
	"github.com/deltegui/phoenix/vars"
	"github.com/gorilla/mux"
)

var mainMapper Mapper = Mapper{
	router: mux.NewRouter(),
}

func Map(mapping Mapping, middlewares ...Middleware) {
	mainMapper.Map(mapping, middlewares...)
}

func MapAll(mappings []Mapping, middlewares ...Middleware) {
	mainMapper.MapAll(mappings, middlewares...)
}

func MapRoot(controllerBuilder injector.Builder) {
	mainMapper.MapRoot(controllerBuilder)
}

func MapGroup(root string, createGroup func(mapper Mapper)) {
	mainMapper.MapGroup(root, createGroup)
}

func Run(listenURL string) {
	printLogo()
	createStaticServer(mainMapper.router)
	showEndpoints(mainMapper.router)
	startListening(listenURL)
}

func createStaticServer(router *mux.Router) {
	if vars.IsStaticServerEnabled() {
		s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
		router.PathPrefix("/static").Handler(s)
		log.Println("Created static server!")
	}
}

func printLogo() {
	if vars.IsLogoFileEnabled() {
		logo, err := ioutil.ReadFile(vars.GetLogoFilename())
		if err != nil {
			log.Fatalf("Cannot read logo file: %s\n", err)
		}
		fmt.Println(string(logo))
	}
	fmt.Print(vars.FormatProjectInfo())
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
