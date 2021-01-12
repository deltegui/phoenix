package phoenix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func CreateTemplates() *template.Template {
	pattern := "./templates/**/*.html"
	templateEngine := template.Must(template.New("html").ParseGlob(pattern))
	log.Printf("Template engine%s\n", templateEngine.DefinedTemplates())
	return templateEngine
}

func formatViewName(view string) string {
	return fmt.Sprintf("%s.html", view)
}

func NewHTMLPresenter(writer http.ResponseWriter, req *http.Request, templates *template.Template, view string) HTMLPresenter {
	realView := formatViewName(view)
	return HTMLPresenter{
		Writer:    writer,
		Request:   req,
		Templates: templates,
		View:      realView,
		ErrorView: realView,
	}
}

func NewHTMLPresenterWithErrView(writer http.ResponseWriter, req *http.Request, templates *template.Template, view string, errView string) HTMLPresenter {
	return HTMLPresenter{
		Writer:    writer,
		Request:   req,
		Templates: templates,
		View:      formatViewName(view),
		ErrorView: formatViewName(errView),
	}
}

type HTMLPresenter struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Templates *template.Template
	View      string
	ErrorView string
}

func (renderer HTMLPresenter) Present(data interface{}) {
	if !renderer.RenderTemplate(renderer.View, data) {
		log.Fatalf("Cannot find view with name: %s\n", renderer.View)
	}
}

func (renderer HTMLPresenter) PresentError(caseError error) {
	renderer.RenderTemplate("error.html", caseError)
}

func (renderer HTMLPresenter) RenderTemplate(view string, data interface{}) bool {
	if err := renderer.Templates.ExecuteTemplate(renderer.Writer, view, data); err != nil {
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
	renderer.Writer.Header().Add("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

func (renderer JSONPresenter) PresentError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Writer.Header().Add("Content-Type", "application/json")
	renderer.Present(caseError)
}
