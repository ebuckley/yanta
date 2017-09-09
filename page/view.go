package page

import (
	"log"
	"net/http"
	"text/template"

	"github.com/ebuckley/yanta/context"
	"github.com/ebuckley/yanta/site"
	"github.com/gorilla/mux"
)

const pageContent = `
<html>
	<head>
		<link href="https://fonts.googleapis.com/css?family=Inconsolata|Nunito+Sans" rel="stylesheet">
		<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js"></script>
		<title>{{.Path}}</title>
		<style>
		 body {
			font-family: 'Nunito Sans', sans-serif;
			color: #3E606F;
		 }
		 h1, h2, h3, h4 {
			font-family: 'Inconsolata', monospace;
		 }
		 a, a:visited {
			 color: #468966;
		 }
		 a:hover {
			 color: #FFB03B;
		 }

		 .content {
			 max-width: 700px;
			 margin: 0 auto 0;
		 }
		</style>
	</head>
	<body>
		<div class='content'>
		 {{.Html}}
		</div>
	</body>
</html>
`

// CreateHandler sets up a handler for serving a single page.
func CreateHandler(page *context.Page) http.HandlerFunc {

	tpl := template.New("page")
	t, err := tpl.Parse(pageContent)
	if err != nil {
		log.Panic("page template couldn't be setup!", page, err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		err = t.Execute(w, page)
		if err != nil {
			log.Fatal("fatal error sending payload", page)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// LookupHandler checks a "path" param in the url
func LookupHandler(s *site.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		page, err := s.FetchPage(path)
		if err != nil {
			log.Println("could not fetch page", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		CreateHandler(page)(w, r)
	}
}
