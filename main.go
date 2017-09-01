package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebuckley/marked/context"
	"github.com/ebuckley/marked/dashboard"
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

	http.HandleFunc("/", dashboard.CreateHandler(pages))

	for _, page := range pages {
		// we need to have a closure over this variable
		thisPage := page
		http.HandleFunc("/"+page.Path, func(w http.ResponseWriter, req *http.Request) {
			_, err = fmt.Fprint(w, thisPage.Html)
			if err != nil {
				log.Fatal("fatal error sending payload", thisPage)
			}
		})

		downloadPdf := "/" + page.Path + "/pdf"
		http.HandleFunc(downloadPdf, pdf.CreatePageDownloader(thisPage))
	}

	log.Println("Started on :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}
