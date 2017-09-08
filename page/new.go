package page

import (
	"html/template"
	"net/http"

	"github.com/ebuckley/yanta/site"
)

// NewPageHandler is a function which sets up the view for creating a new page
func NewPageHandler(s *site.Site) http.HandlerFunc {
	t := template.Must(template.New("new").Parse(upsertPageContent))
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
