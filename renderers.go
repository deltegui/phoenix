package phoenix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func parseTemplate(layout, view string) *template.Template {
	templateRoot := "web/templates"
	path := fmt.Sprintf("%s%s", templateRoot, view)
	layoutPath := fmt.Sprintf("%s%s", templateRoot, layout)
	if layout == "" {
		return template.Must(template.ParseFiles(path))
	} else {
		return template.Must(template.ParseFiles(layoutPath, path))
	}
}

type RequestMapper func(req *http.Request) interface{}

type ViewConfig struct {
	Layout, View, Name string
}

func RenderView(conf ViewConfig, mapper RequestMapper) http.HandlerFunc {
	tmpl := parseTemplate(conf.Layout, conf.View)
	return func(w http.ResponseWriter, req *http.Request) {
		tmpl.ExecuteTemplate(w, conf.Name, mapper(req))
	}
}

type HTMLRenderer struct {
	view     string
	template *template.Template
}

func NewHTMLRenderer(conf ViewConfig) HTMLRenderer {
	return HTMLRenderer{
		view:     conf.Name,
		template: parseTemplate(conf.Layout, conf.View),
	}
}

func (renderer HTMLRenderer) execute(w http.ResponseWriter, viewmodel interface{}) {
	renderer.template.ExecuteTemplate(w, renderer.view, viewmodel)
}

func (renderer HTMLRenderer) Render(w http.ResponseWriter, data interface{}) {
	renderer.execute(w, data)
}

// JSONPresenter is a presenter that renders your data in JSON format.
type JSONPresenter struct {
	Writer http.ResponseWriter
}

// NewJSONPresenter creates a presenter that renders your data in JSON format.
func NewJSONPresenter(writer http.ResponseWriter) JSONPresenter {
	return JSONPresenter{writer}
}

// Present data in JSON format
func (renderer JSONPresenter) Present(data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		return
	}
	renderer.Writer.Header().Set("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

// PresentError renders a JSON with your error
func (renderer JSONPresenter) PresentError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Present(caseError)
}
