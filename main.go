package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ebuckley/marked/dashboard"
	"github.com/ebuckley/marked/page"
	"github.com/ebuckley/marked/pdf"
	"github.com/ebuckley/marked/site"
	"github.com/gorilla/mux"
)

var (
	bufferLength = 1024
)

func makePageMatcher(s string) mux.MatcherFunc {
	re := regexp.MustCompile("/" + s + "/.*$")
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		match := re.MatchString(r.URL.Path)
		// TODO handle error from regex
		if match {
			path := strings.TrimPrefix(r.URL.Path, "/"+s+"/")
			rm.Vars = make(map[string]string)
			rm.Vars["path"] = path
		}
		return match
	}
}

func main() {
	var appDir string
	flag.StringVar(&appDir, "dir", "./docs", "site directory")
	flag.Parse()

	s := site.New(site.SitePath(appDir))

	r := mux.NewRouter()

	r.HandleFunc("/", dashboard.CreateHandler(s))

	r.Methods("GET").
		MatcherFunc(makePageMatcher("page")).
		HandlerFunc(page.LookupHandler(s))

	r.Methods("POST").
		MatcherFunc(makePageMatcher("page")).
		HandlerFunc(page.UpsertPageHandler(s))

	r.MatcherFunc(makePageMatcher("pdf")).
		HandlerFunc(pdf.LookupHandler(s))

	r.Methods("GET").
		MatcherFunc(makePageMatcher("edit")).
		HandlerFunc(page.ViewUpsert(s))

	http.Handle("/", r)

	log.Println("Started on :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}
