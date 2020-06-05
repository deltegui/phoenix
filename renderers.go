package phoenix

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var templateEngine *template.Template = nil

func getTemplateEngine() *template.Template {
	if templateEngine != nil {
		return templateEngine
	}
	templateEngine = template.Must(template.New("html").ParseGlob("./templates/*.html"))
	log.Printf("Template engine %s\n", templateEngine.DefinedTemplates())
	return templateEngine
}

type HTMLPresenter struct {
	Writer http.ResponseWriter
	View   string
}

func NewHTMLPresenter(writer http.ResponseWriter, view string) HTMLPresenter {
	return HTMLPresenter{
		Writer: writer,
		View:   view,
	}
}

func (renderer HTMLPresenter) Present(data interface{}) {
	if !renderer.renderTemplate(renderer.View, data) {
		log.Fatalf("Cannot find view with name: %s\n", renderer.View)
	}
}

func (renderer HTMLPresenter) PresentError(caseError error) {
	renderer.renderTemplate("error.html", caseError)
}

func (renderer HTMLPresenter) renderTemplate(view string, data interface{}) bool {
	if err := getTemplateEngine().ExecuteTemplate(renderer.Writer, view, data); err != nil {
		log.Print("Error during rendering template: ")
		log.Println(err)
		return false
	}
	return true
}

type JSONPresenter struct {
	Writer http.ResponseWriter
}

func NewJSONPresenter(writer http.ResponseWriter) JSONPresenter {
	return JSONPresenter{writer}
}

func (renderer JSONPresenter) Present(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	renderer.Writer.Header().Set("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

func (renderer JSONPresenter) PresentError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Present(caseError)
}
