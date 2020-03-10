package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clevergo/views"
)

var manager *views.Manager

func init() {
	viewsPath := "views" // views path.
	// options
	opts := []views.Option{
		// views.Suffix(".tmpl"), // template suffix, default to .tmpl.
		// views.Delims("{{", "}}"), // template delimiters, default to "{{" and "}}".
		views.DefaultLayout("main", "head", "header", "footer"),
		views.LayoutsDir("layouts"),   // layout directory, relatived to views path.
		views.PartialsDir("partials"), // partials layout, relatived to layouts directory.
		// global function map for all templates.
		views.FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
		views.Cache(false), // disabled caching for developing.
	}
	manager = views.New(viewsPath, opts...)
	// add a new layout.
	manager.AddLayout("page", "head")

	// regiters before render listener.
	manager.RegisterOnBeforeRender(func(event *views.BeforeRenderEvent) error {
		return nil
	})

	// regiters after render listener.
	manager.RegisterOnAfterRender(func(event *views.AfterRenderEvent) error {
		return nil
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	manager.Render(w, "site/index", views.Context{
		"title": "home",
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	manager.RenderLayout(w, "page", "user/login", nil)
}

func partial(w http.ResponseWriter, r *http.Request) {
	manager.RenderPartial(w, "site/partial", views.Context{
		"title": "partial",
	})
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/partial", partial)
	http.ListenAndServe(":1234", http.DefaultServeMux)
}
