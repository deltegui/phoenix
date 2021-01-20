package phoenix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

// CreateTemplates using the pattern './templates/**/*.html'.
func CreateTemplates() {
	pattern := "./templates/**/*.html"
	templates = template.Must(template.New("html").ParseGlob(pattern))
	log.Printf("Template engine%s\n", templates.DefinedTemplates())
}

func formatViewName(view string) string {
	return fmt.Sprintf("%s.html", view)
}

// RenderTemplate using a response writer, a view name and data to pass to the view.
// Returns false if there is an error duing template rendering. Returns true if
// everything is ok.
// View names are expected to have .html extension, so it's no needed to pass the name
// with the extension. For example, if your view is in ./templates/users/user_index.html
// to render it you have to call this function this way:
//
// 		phoenix.RenderTemplate(w, "user_index", nil)
//
// The data parameter can be nil if you dont want to pass anything to your view.
func RenderTemplate(w http.ResponseWriter, view string, data interface{}) bool {
	if err := templates.ExecuteTemplate(w, formatViewName(view), data); err != nil {
		log.Print("Error during rendering template: ")
		log.Println(err)
		return false
	}
	return true
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
	renderer.Writer.Header().Add("Content-Type", "application/json")
	renderer.Writer.Write(response)
}

// PresentError renders a JSON with your error
func (renderer JSONPresenter) PresentError(caseError error) {
	renderer.Writer.WriteHeader(http.StatusBadRequest)
	renderer.Present(caseError)
}
