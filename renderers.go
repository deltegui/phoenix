package phoenix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/deltegui/phoenix/vars"
)

type RenderMetadata struct {
	ViewName string
}

type Renderer interface {
	Render(data interface{})
	RenderWithMeta(data interface{}, metadata RenderMetadata)
	RenderError(data error)
}

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

func (renderer HTMLRenderer) Render(data interface{}) {
	found := false
	for skip := 3; skip > 1 && !found; skip-- {
		view := viewNameFromCallerSkipping(skip)
		log.Printf("Extracted view name: %s\n", view)
		if renderer.renderTemplate(view, data) {
			found = true
		}
	}
	if !found {
		log.Fatalln("Cannot find your view")
	}
}

func (renderer HTMLRenderer) RenderWithMeta(data interface{}, metadata RenderMetadata) {
	if !renderer.renderTemplate(metadata.ViewName, data) {
		log.Fatalf("Cannot find view with name: %s\n", metadata.ViewName)
	}
}

func (renderer HTMLRenderer) RenderError(caseError error) {
	renderer.RenderWithMeta(caseError, RenderMetadata{"error.html"})
}

func (renderer HTMLRenderer) renderTemplate(view string, data interface{}) bool {
	if err := getTemplateEngine().ExecuteTemplate(renderer.Writer, view, data); err != nil {
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

func (renderer JSONRenderer) RenderWithMeta(data interface{}, metadata RenderMetadata) {
	renderer.Render(data)
}

func (renderer JSONRenderer) RenderError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Render(caseError)
}
