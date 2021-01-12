package phoenix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func CreateTemplates() {
	pattern := "./templates/**/*.html"
	templates = template.Must(template.New("html").ParseGlob(pattern))
	log.Printf("Template engine%s\n", templates.DefinedTemplates())
}

func formatViewName(view string) string {
	return fmt.Sprintf("%s.html", view)
}

func RenderTemplate(w http.ResponseWriter, view string, data interface{}) bool {
	if err := templates.ExecuteTemplate(w, formatViewName(view), data); err != nil {
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
