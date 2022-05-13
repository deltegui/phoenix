package phoenix

import (
	"encoding/json"
	"log"
	"net/http"
)

type Present func(data interface{}, errs []error)

func JSONPresenter(w http.ResponseWriter, req *http.Request) Present {
	render := func(data interface{}) {
		response, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshaling data: ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
	return func(data interface{}, errs []error) {
		if errs != nil {
			w.WriteHeader(http.StatusBadRequest)
			render(errs)
			return
		}
		render(data)
	}
}
