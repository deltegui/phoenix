package phoenix

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/deltegui/phoenix/vars"
)

var templateEngine *template.Template = nil

func getTemplateEngine() *template.Template {
	if templateEngine != nil {
		return templateEngine
	}
	if vars.IsTemplatesEnabled() {
		templateEngine = template.Must(template.New("html").ParseGlob("./templates/*/*.html"))
		log.Printf("Template engine %s\n", templateEngine.DefinedTemplates())
	} else {
		log.Fatalln("Trying to get template engine when its disabled")
	}
	return templateEngine
}

type HTMLRenderer struct {
	Writer http.ResponseWriter
}

func (renderer HTMLRenderer) Render(view string, data interface{}) {
	if !renderer.renderTemplate(view, data) {
		log.Fatalf("Cannot find view with name: %s\n", view)
	}
}

func (renderer HTMLRenderer) RenderError(caseError error) {
	renderer.Render("error.html", caseError)
}

func (renderer HTMLRenderer) renderTemplate(view string, data interface{}) bool {
	if err := getTemplateEngine().ExecuteTemplate(renderer.Writer, view, data); err != nil {
		log.Print("Error during rendering template: ")
		log.Println(err)
		return false
	}
	return true
}

type JSONRenderer struct {
	Writer http.ResponseWriter
}

func (renderer JSONRenderer) Render(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	renderer.Writer.Header().Set("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

func (renderer JSONRenderer) RenderError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Render(caseError)
}
