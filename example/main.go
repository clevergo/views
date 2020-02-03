package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clevergo/views"
)

var view *views.View

func init() {
	// options
	opts := []views.Option{
		// views.Suffix(".tmpl"), // template suffix, default to .tmpl.
		// views.Delims("{{", "}}"), // template delimiters, default to "{{" and "}}".
		views.Theme("default"),
		views.Layouts("layouts/main", "layouts/header", "layouts/footer"),
		// global function map for all templates.
		views.FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
		// views.Cache(true), // disabled it for developing.
	}
	view = views.New("./themes", opts...)
}

func home(w http.ResponseWriter, r *http.Request) {
	err := view.Render(w, "site/index", map[string]interface{}{
		"title": "home",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func partial(w http.ResponseWriter, r *http.Request) {
	err := view.RenderPartial(w, "site/partial", map[string]interface{}{
		"title": "standalone",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/partial", partial)
	http.ListenAndServe(":1234", http.DefaultServeMux)
}
