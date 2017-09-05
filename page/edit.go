package page

import (
	"log"
	"net/http"
	"text/template"

	"github.com/ebuckley/yanta/site"
	"github.com/gorilla/mux"
)

const editPageContent = `
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
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">
		<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>
		<script>
		document.addEventListener("DOMContentLoaded", function(event) {
			console.log("DOM fully loaded and parsed");
			var simplemde = new SimpleMDE({
				autofocus: true,
				spellChecker: false
			});
		});

		</script>
	</head>
	<body>
			<div class=".content">
				<form action="/page/{{.Path}}" name="upsertForm" method="post">
					<div>
						<input class="saver" type="submit" value="Save">
					</div>
					<textarea name="content">{{.Content}}</textarea>
				</form>
			</div>
	</body>
</html>
`

// ViewUpsert is for looking at the edit view :O
func ViewUpsert(s *site.Site) http.HandlerFunc {

	tpl := template.New("upsert")
	t, err := tpl.Parse(editPageContent)
	if err != nil {
		log.Panic("page template couldn't be setup!", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]

		page, err := s.FetchPage(path)
		if err != nil {
			log.Println("edit view fail with fetching path", path)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = t.Execute(w, page)
		if err != nil {
			log.Println("edit view fail with template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
