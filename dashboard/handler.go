package dashboard

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ebuckley/yanta/site"
)

const pageContent = `
<html>
	<head>
		<link href="https://fonts.googleapis.com/css?family=Inconsolata|Nunito+Sans" rel="stylesheet">
		<title>Get things done!</title>
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
		 .btn-primary {
			 border: 1px solid #a2a2a2;
			 border-radius: 3px;
			 padding: 4px;
			 text-transform: uppercase;
			 color: #3E606F;
		 }
		 .nav {
			 margin-bottom: 16px;
		 }
		</style>
		<meta http-equiv="refresh" content="2">
	</head>
	<body>
		<h1>{{.Config.ApplicationName}}</h1>
		<div class='nav'>
			<a class='btn-primary' href="/new">new</a>
			{{if .CanPublish}}
				<a class='btn-primary' href="/publish">Publish</a>
				<a class='btn-primary' href="/pull">Pull</a>
			{{end}}
		</div>
		<table>
		{{range .Pages}}
			<tr>
				<td>
					<a href="/page/{{.Path}}" target="_blank" >
						{{.Path}}
					</a>
				</td>
				<td>
					<a href="/pdf/{{.Path}}" target="_blank" >
					pdf
					</a>
				</td>
				<td>
				  <a href="/edit/{{.Path}}" target="_blank">edit</a>
				</td>
			</tr>
		{{end}}
		</table>
	</body>
</html>
`

// CreateHandler sets up a dashboard page handler
func CreateHandler(s *site.Site) http.HandlerFunc {
	tpl := template.New("index")
	t, err := tpl.Parse(pageContent)
	if err != nil {
		log.Panic("Template setup not possible", err)
	}
	//
	return func(w http.ResponseWriter, r *http.Request) {
		err = s.FetchPages()
		if err != nil {
			log.Fatal("could not fetch pages", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, s)
		if err != nil {
			log.Fatal("could not execute template", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
