package phoenix

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

type App struct {
	Mapper
	config   *PhoenixConfig
	Injector *Injector
}

func NewApp() App {
	injector := NewInjector()
	return App{
		Mapper: Mapper{
			router:   mux.NewRouter(),
			injector: injector,
		},
		config: &PhoenixConfig{
			projectName:        "phoenix",
			projectVersion:     "0.1.0",
			enableStaticServer: false,
			logoFile:           "",
			enableSessions:     false,
			onStop:             func() {},
			tlsCertFile:        "",
			tlsKeyFile:         "",
			domains:            nil,
		},
		Injector: injector,
	}
}

func (app App) Configure() *PhoenixConfig {
	return app.config
}

func (app App) Run(listenURL string) {
	app.printLogo()
	app.createStaticServer()
	app.showEndpoints()
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
	server := &http.Server{
		Addr:    address,
		Handler: app.Mapper.router,
	}
	go app.startServer(server)
	app.waitAndStopServer(server)
}

func (app App) startServer(server *http.Server) {
	if app.config.isAutoHTTPSEnabled() {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("golang-autocert"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(app.config.getAutoHTTPSDomains()...),
		}
		server.TLSConfig = m.TLSConfig()
	}
	log.Println("Listening on address: ", server.Addr)
	log.Println("You are ready to GO!")
	var err error
	if app.config.isAutoHTTPSEnabled() {
		err = server.ListenAndServeTLS("", "")
	} else if app.config.isHTTPSEnabled() {
		certFile, keyFile := app.config.getHTTPSCertKeyFiles()
		err = server.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Fatalln("Error while listening: ", err)
	}
}

func (app App) waitAndStopServer(server *http.Server) {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		app.config.getStopHandler()()
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Phoenix shutdown failed:%+v", err)
	}

	log.Print("Phoenix exited properly")
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
