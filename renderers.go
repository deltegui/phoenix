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
	if layout == "" {
		return template.Must(template.ParseFiles(path))
	} else {
		return template.Must(template.ParseFiles(layout, path))
	}
}

type CreateViewModel func(data interface{}) interface{}
type CreateErrorViewModel func(err error) interface{}

type HTMLPresenter struct {
	view                 string
	createViewModel      CreateViewModel
	createErrorViewModel CreateErrorViewModel
	template             *template.Template
	w                    http.ResponseWriter
}

func NewHTMLPresenter(view string, cvm CreateViewModel, cevm CreateErrorViewModel) HTMLPresenter {
	return HTMLPresenter{
		view:                 view,
		createViewModel:      cvm,
		createErrorViewModel: cevm,
		template:             parseTemplate("", view),
	}
}

func NewHTMLPresenterWithLayout(layout, view string, cvm CreateViewModel, cevm CreateErrorViewModel) HTMLPresenter {
	return HTMLPresenter{
		view:                 view,
		createViewModel:      cvm,
		createErrorViewModel: cevm,
		template:             parseTemplate(layout, view),
	}
}

func (presenter *HTMLPresenter) Use(w http.ResponseWriter) {
	presenter.w = w
}

func (presenter HTMLPresenter) execute(viewmodel interface{}) {
	presenter.template.ExecuteTemplate(presenter.w, presenter.view, viewmodel)
}

func (presenter HTMLPresenter) Present(data interface{}) {
	presenter.execute(presenter.createViewModel(data))
}

func (presenter HTMLPresenter) PresentError(err error) {
	presenter.execute(presenter.createErrorViewModel(err))
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
