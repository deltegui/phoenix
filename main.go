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
