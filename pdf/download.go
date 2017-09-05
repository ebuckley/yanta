package pdf

import (
	"log"
	"net/http"

	"github.com/ebuckley/yanta/context"
	"github.com/ebuckley/yanta/site"
	"github.com/gorilla/mux"
)

// CreatePageDownloader creates a handlerfunction for downloading a page
func CreatePageDownloader(cfg *site.Config, page *context.Page) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", `inline`)
		data, err := Download(cfg, page)
		if err != nil {
			log.Fatal("failed to download ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = w.Write(data)
		if err != nil {
			log.Fatal("failed to write content", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// LookupHandler is a function for finding the page for a given request
func LookupHandler(s *site.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		p, err := s.FetchPage(path)
		if err != nil {
			log.Println("could not fetch page!", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		CreatePageDownloader(s.Config, p)(w, r)
	}
}
