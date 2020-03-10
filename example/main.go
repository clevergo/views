package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/views/v2"
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

func home(ctx *clevergo.Context) error {
	return manager.Render(ctx.Response, "site/index", views.Context{
		"title": "home",
	})
}

func login(ctx *clevergo.Context) error {
	return manager.RenderLayout(ctx.Response, "page", "user/login", nil)
}

func partial(ctx *clevergo.Context) error {
	return manager.RenderPartial(ctx.Response, "site/partial", views.Context{
		"title": "partial",
	})
}

func main() {
	router := clevergo.NewRouter()
	router.Get("/", home)
	router.Get("/login", login)
	router.Get("/partial", partial)
	http.ListenAndServe(":1234", router)
}
