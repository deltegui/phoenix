package locomotive

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/deltegui/locomotive/vars"
	"github.com/gorilla/mux"
)

func Run(listenURL string) {
	printLogo()
	createStaticServer(router)
	showEndpoints(router)
	startListening(listenURL)
}

func createStaticServer(router *mux.Router) {
	if vars.EnableStaticServer {
		s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
		router.PathPrefix("/static").Handler(s)
		log.Println("Created static server!")
	}
}

func printLogo() {
	if vars.LogoFile != "" {
		logo, err := ioutil.ReadFile(vars.LogoFile)
		if err != nil {
			log.Fatalf("Cannot read logo file: %s\n", err)
		}
		fmt.Println(string(logo))
	}
	fmt.Printf("%s v%s\n", vars.Name, vars.Version)
}
