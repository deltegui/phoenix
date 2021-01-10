package phoenix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templateEngine *template.Template = createTemplates()

func createTemplates() *template.Template {
	pattern := "./templates/**/*.html"
	templateEngine := template.Must(template.New("html").ParseGlob(pattern))
	log.Printf("Template engine%s\n", templateEngine.DefinedTemplates())
	return templateEngine
}

func formatViewName(view string) string {
	return fmt.Sprintf("%s.html", view)
}

func BuildPresenter(writer http.ResponseWriter, req *http.Request, view string) HTMLPresenter {
	realView := formatViewName(view)
	return HTMLPresenter{
		Writer:    writer,
		Request:   req,
		View:      realView,
		ErrorView: realView,
	}
}

func BuildPresenterWithErr(writer http.ResponseWriter, req *http.Request, view string, errView string) HTMLPresenter {
	return HTMLPresenter{
		Writer:    writer,
		Request:   req,
		View:      formatViewName(view),
		ErrorView: formatViewName(errView),
	}
}

type HTMLPresenter struct {
	Writer    http.ResponseWriter
	Request   *http.Request
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
	if err := templateEngine.ExecuteTemplate(renderer.Writer, view, data); err != nil {
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
