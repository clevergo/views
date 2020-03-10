package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clevergo/views"
)

var manager *views.Manager

func init() {
	// options
	opts := []views.Option{
		// views.Suffix(".tmpl"), // template suffix, default to .tmpl.
		// views.Delims("{{", "}}"), // template delimiters, default to "{{" and "}}".
		views.DefaultLayout("main", "head", "header", "footer"),
		// global function map for all templates.
		views.FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
		// views.Cache(true), // disabled it for developing.
	}
	manager = views.New("views", opts...)
	manager.AddLayout("page", "head")
}

func home(w http.ResponseWriter, r *http.Request) {
	manager.Render(w, "site/index", views.Context{
		"title": "home",
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	manager.RenderLayout(w, "page", "site/login", views.Context{
		"title": "standalone",
	})
}

func partial(w http.ResponseWriter, r *http.Request) {
	manager.RenderPartial(w, "site/partial", views.Context{
		"title": "standalone",
	})
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/partial", partial)
	http.ListenAndServe(":1234", http.DefaultServeMux)
}
