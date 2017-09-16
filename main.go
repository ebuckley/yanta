package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ebuckley/yanta/dashboard"
	"github.com/ebuckley/yanta/page"
	"github.com/ebuckley/yanta/pdf"
	"github.com/ebuckley/yanta/site"
	"github.com/ebuckley/yanta/sync"
	"github.com/gorilla/mux"
)

var (
	bufferLength = 1024
)

func makePageMatcher(s string) mux.MatcherFunc {
	re := regexp.MustCompile("/" + s + "/.*$")
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		match := re.MatchString(r.URL.Path)
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
	var configPath string
	flag.StringVar(&appDir, "dir", "./docs", "site directory")
	flag.Parse()

	flag.StringVar(&configPath, "config", "yanta.json", "The path of your yanta config file")

	cfg, err := site.DecodeConfig(configPath)
	if err != nil {
		log.Println("Could not load config", err)
	}
	var s *site.Site
	if cfg != nil {
		log.Println("Creating site from config", cfg)
		s = site.New(site.FromConfig(cfg))
	} else {
		s = site.New(site.SitePath(appDir))
	}

	r := mux.NewRouter()

	r.HandleFunc("/", dashboard.CreateHandler(s))

	r.Methods("GET").
		MatcherFunc(makePageMatcher("page")).
		HandlerFunc(page.LookupHandler(s))

	r.Methods("POST").
		MatcherFunc(makePageMatcher("page")).
		HandlerFunc(page.UpsertPageHandler(s))

	r.Methods("GET").
		MatcherFunc(makePageMatcher("edit")).
		HandlerFunc(page.ViewUpsert(s))

	r.Methods("GET").
		Path("/new").
		HandlerFunc(page.NewPageHandler(s))

	r.Methods("POST").
		Path("/new").
		HandlerFunc(page.UpsertPageHandler(s))

	r.MatcherFunc(makePageMatcher("pdf")).
		HandlerFunc(pdf.LookupHandler(s))

	r.MatcherFunc(makePageMatcher("history")).
		HandlerFunc(sync.HistoryHandler(s))

	r.HandleFunc("/publish", sync.PublishHandler(s))
	r.HandleFunc("/pull", sync.PullHandler(s))

	r.Methods("GET").
		Path("/config.json").
		HandlerFunc(site.ServeConfig(s))

	r.Methods("POST").
		Path("/config.json").
		HandlerFunc(site.UpdateConfig(s))

	r.Methods("OPTIONS").
		Path("/config.json").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})

	http.Handle("/", middleware(r))

	log.Println("Started on :1337")
	log.Fatal(http.ListenAndServe(":1337", nil))
}

func middleware(h http.Handler) http.Handler {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Method", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		h.ServeHTTP(w, r)
	})
}
