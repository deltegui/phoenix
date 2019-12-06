package locomotive

import (
	"encoding/json"
	"fmt"
	"html/template"
	"locomotive/vars"
	"log"
	"net/http"
	"runtime"
	"strings"
)

var templateEngine *template.Template = nil

func getTemplateEngine() *template.Template {
	if templateEngine != nil {
		return templateEngine
	}
	if vars.EnableTemplates {
		templateEngine = template.Must(template.New("html").ParseGlob("./templates/*/*.html"))
		log.Printf("Template engine %s\n", templateEngine.DefinedTemplates())
	} else {
		log.Fatalln("Trying to get template engine when its disabled")
	}
	return templateEngine
}

type htmlPresenter struct {
	writer http.ResponseWriter
}

//Present iformation in HTML
func (presenter htmlPresenter) Present(data interface{}) {
	found := false
	for skip := 3; skip > 1 && !found; skip-- {
		view := viewNameFromCallerSkipping(skip)
		log.Printf("Extracted view name: %s\n", view)
		if presenter.render(view, data) {
			found = true
		}
	}
	if !found {
		log.Fatalln("Cannot find your view")
	}
}

//PresentError in HTML
func (presenter htmlPresenter) PresentError(caseError error) {
	presenter.render("error.html", caseError)
}

func (presenter htmlPresenter) render(view string, data interface{}) bool {
	if err := getTemplateEngine().ExecuteTemplate(presenter.writer, view, data); err != nil {
		log.Print("Error during rendering template: ")
		log.Println(err)
		return false
	}
	return true
}

func viewNameFromCallerSkipping(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		log.Fatalf("HTML presenter: Cannot obtain runtime caller!")
	}
	functionCaller := runtime.FuncForPC(pc)
	fullName := functionCaller.Name()
	tokens := strings.Split(fullName, ".")
	return fmt.Sprintf("%s.html", strings.ToLower(tokens[len(tokens)-1]))
}

type jsonPresenter struct {
	writer http.ResponseWriter
}

//Present data in JSON format
func (presenter jsonPresenter) Present(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	presenter.writer.Header().Set("Content-Type", "application/json")
	presenter.writer.Write(response)
}

//PresentError in JSON format
func (presenter jsonPresenter) PresentError(caseError error) {
	presenter.writer.WriteHeader(http.StatusBadRequest)
	presenter.Present(caseError)
}
