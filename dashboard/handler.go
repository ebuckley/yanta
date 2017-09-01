package dashboard

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ebuckley/marked/context"
)

const pageContent = `
<html>
	<head>
		<title>Get things done!</title>
	</head>
	<body>
		<table>
		{{range .}}
			<tr>
				<td>
					<a href="/{{.Path}}" target="_blank" >
						{{.Path}}
					</a>
				</td>
				<td>
					<a href="/{{.Path}}/pdf" target="_blank" >
					pdf
					</a>
				</td>
			</tr>
		{{end}}
		</table>
	</body>
</html>
`

// CreateHandler sets up a dashboard page handler
func CreateHandler(p []*context.Page) http.HandlerFunc {
	tpl := template.New("index")
	t, err := tpl.Parse(pageContent)
	if err != nil {
		log.Panic("Template setup not possible", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err = t.Execute(w, p)
		if err != nil {
			log.Fatal("could not execute template", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
