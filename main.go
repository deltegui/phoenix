package phoenix

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type App struct {
	Mapper   Mapper
	config   PhoenixConfig
	Injector *Injector
}

func NewApp() App {
	injector := NewInjector()
	return App{
		Mapper: Mapper{
			router:   mux.NewRouter(),
			injector: injector,
		},
		config: PhoenixConfig{
			projectName:        "phoenix",
			projectVersion:     "0.1.0",
			enableStaticServer: false,
			enableTemplates:    false,
			logoFile:           "",
		},
		Injector: injector,
	}
}

func (app App) Configure() PhoenixConfig {
	return app.config
}

func (app App) Run(listenURL string) {
	app.printLogo()
	app.createStaticServer(mainMapper.router)
	app.showEndpoints(mainMapper.router)
	app.startListening(listenURL)
}

func (app App) createStaticServer() {
	if app.config.isStaticServerEnabled() {
		s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
		app.Mapper.router.PathPrefix("/static").Handler(s)
		log.Println("Created static server!")
	}
}

func (app App) printLogo() {
	if app.config.isLogoFileEnabled() {
		logo, err := ioutil.ReadFile(app.config.getLogoFilename())
		if err != nil {
			log.Fatalf("Cannot read logo file: %s\n", err)
		}
		fmt.Println(string(logo))
	}
	fmt.Print(app.config.formatProjectInfo())
}

func (app App) startListening(address string) {
	log.Println("Listening on address: ", address)
	log.Println("You are ready to GO!")
	if err := http.ListenAndServe(address, mainMapper.router); err != nil {
		log.Fatalln("Error while listening: ", err)
	}
}

func (app App) showEndpoints() {
	stringEndsWith := func(target, end string) bool {
		return strings.LastIndex(target, end) == len(target)-1
	}

	app.Mapper.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
