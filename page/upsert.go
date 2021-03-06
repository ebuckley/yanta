package page

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/ebuckley/yanta/site"
	"github.com/ebuckley/yanta/sync"
)

type upsertRequest struct {
	PageName string `schema:"pagename"`
	Content  string `schema:"content"`
}

// UpsertPageHandler is a function which returns a handler for updating a page!
func UpsertPageHandler(s *site.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]

		err := r.ParseForm()

		if err != nil {
			log.Println("Could not parse form", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		upsertReq := new(upsertRequest)

		decoder := schema.NewDecoder()
		// r.PostForm is a map of our POST form values
		err = decoder.Decode(upsertReq, r.PostForm)

		if err != nil {
			// Handle error
			log.Println("Could not handle request..", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if upsertReq.PageName != "" {
			path = filepath.Join(s.Config.SitePath, upsertReq.PageName)
		}

		// add a file!
		p, err := s.UpsertPage(path, upsertReq.Content)
		if err != nil {
			log.Println("big problem updating the page", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hst, _ := os.Hostname()
		msg := fmt.Sprint("automatic commit by ", hst)
		out, err := sync.Commit(p, msg)
		if err != nil {
			log.Println("saved file, but could not commit change")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("success:", out)

		http.Redirect(w, r, "/page/"+p.Path, http.StatusFound)
	}
}
