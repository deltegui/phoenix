package locomotive

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/deltegui/locomotive/vars"
)

type PresenterMetadata struct {
	ViewName string
}

type Presenter interface {
	Present(data interface{})
	Present(data interface{}, metadata PresenterMetadata)
	PresentError(data error)
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

type HTMLPresenter struct {
	Writer http.ResponseWriter
}

func (presenter HTMLPresenter) Present(data interface{}) {
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

func (presenter HTMLPresenter) Present(data interface{}, metadata PresenterMetadata) {
	if !presenter.render(metadata.ViewName, data) {
		log.Fatalf("Cannot find view with name: %s\n", metadata.ViewName)
	}
}

func (presenter HTMLPresenter) PresentError(caseError error) {
	presenter.Present(caseError, PresenterMetadata{"error.html"})
}

func (presenter HTMLPresenter) render(view string, data interface{}) bool {
	if err := getTemplateEngine().ExecuteTemplate(presenter.Writer, view, data); err != nil {
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

type JSONPresenter struct {
	Writer http.ResponseWriter
}

func (presenter JSONPresenter) Present(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	presenter.Writer.Header().Set("Content-Type", "application/json")
	presenter.Writer.Write(response)
}

func (presenter HTMLPresenter) Present(data interface{}, metadata PresenterMetadata) {
	presenter.Present(data)
}

func (presenter JSONPresenter) PresentError(caseError error) {
	presenter.Writer.WriteHeader(http.StatusBadRequest)
	presenter.Present(caseError)
}
