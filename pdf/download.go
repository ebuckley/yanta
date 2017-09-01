package pdf

import (
	"log"
	"net/http"

	"github.com/ebuckley/marked/context"
)

// CreatePageDownloader creates a handlerfunction for downloading a page
func CreatePageDownloader(page *context.Page) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", `inline`)
		data, err := Download(page)
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
