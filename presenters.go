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

type Presenter interface {
	Present(data interface{})
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

type HTMLMetadata struct {
	ViewName string
}

type HTMLPresenter struct {
	Writer http.ResponseWriter
}

func (presenter HTMLPresenter) Present(data interface{}) {
	metadata := data.(HTMLMetadata)
	if !presenter.render(metadata.ViewName, data) {
		log.Fatalln("Cannot find your view")
	}
}

func (presenter HTMLPresenter) PresentError(caseError error) {
	presenter.render("error.html", caseError)
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

func (presenter JSONPresenter) PresentError(caseError error) {
	presenter.Writer.WriteHeader(http.StatusBadRequest)
	presenter.Present(caseError)
}
