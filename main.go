package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebuckley/marked/context"
	"github.com/ebuckley/marked/pdf"
)

var (
	bufferLength = 1024
)

const sitePath = "./site"

func main() {
	pages, err := context.GetHTMLs(sitePath)
	if err != nil {
		log.Fatal("It all broke!")
	}
	fmt.Println("marked loaded these files.\n", pages)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello index!")
	})

	for _, page := range pages {
		log.Println("setup for", "/"+page.Path)
		http.HandleFunc("/"+page.Path, func(w http.ResponseWriter, req *http.Request) {
			log.Println("handle for", "/"+page.Path)
			_, err = fmt.Fprint(w, page.Html)
			if err != nil {
				log.Fatal("fatal error sending payload", page)
			}
		})

		downloadPdf := "/" + page.Path + "/pdf"
		http.HandleFunc(downloadPdf, func(w http.ResponseWriter, req *http.Request) {
			data, err := pdf.Download(page)
			if err != nil {
				log.Fatal("failed to download ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			_, err = w.Write(data)
			if err != nil {
				log.Fatal("failed to write content", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}

	log.Println("Started on :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}
