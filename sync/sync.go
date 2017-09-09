package sync

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ebuckley/yanta/context"
	"github.com/ebuckley/yanta/site"
	"github.com/gorilla/mux"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func handleErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println("sync error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createCtx(c object.CommitIter, p *context.Page) (*historyCtx, error) {
	ctx := new(historyCtx)
	ctx.Page = p
	ctx.Commits = make([]*object.Commit, 0)

	err := c.ForEach(func(c *object.Commit) error {
		ctx.Commits = append(ctx.Commits, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

type historyCtx struct {
	Page    *context.Page
	Commits []*object.Commit
}

const historyTpl = `
<html>
	<head>
		<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
	</head>
	<body>
		<table class="pure-table">
			<thead>
				<tr>
						<th>#</th>
						<th>Name</th>
						<th>Message</th>
						<th></th>
				</tr>
			</thead>
			<tbody>
			{{$page := .Page}}
			{{range .Commits}}
				<tr>
					<td>{{.Author.When.Format "2006-01-02 15:04:05"}}</td>
					<td>{{.Author.Name}}</td>
					<td>{{.Message}}</td>
					<td><a href="/page/{{$page.Path}}?commit={{.Hash}}">checkout</a></td>
				</tr>
			{{end}}
			</tdody>
		</table>
	</body>
</html>
`

// HistoryHandler is a handler which returns past instance
func HistoryHandler(s *site.Site) http.HandlerFunc {
	repo, err := git.PlainOpen(s.Config.SitePath)
	if err != nil {
		log.Panic("err opening git", err)
	}

	tpl := template.New("history")
	t, err := tpl.Parse(historyTpl)
	if err != nil {
		log.Panic("not able to parse template")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]

		page, err := s.FetchPage(path)
		handleErr(w, err)

		h, err := repo.Head()

		handleErr(w, err)

		ctx, err := repo.Log(&git.LogOptions{From: h.Hash()})
		handleErr(w, err)

		c, err := createCtx(ctx, page)
		handleErr(w, err)

		err = t.Execute(w, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
