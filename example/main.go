package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/clevergo/clevergo"
	_ "github.com/clevergo/views/example/static/assets"
	_ "github.com/clevergo/views/example/static/views"
	"github.com/clevergo/views/v2"
	"github.com/gobuffalo/packr/v2"
	"github.com/rakyll/statik/fs"
)

var (
	flagFS   = flag.String("fs", "", "file system: packr or statik")
	flagAddr = flag.String("addr", ":8080", "listen address")
	manager  *views.Manager
	viewsFS  http.FileSystem
	assetsFS http.FileSystem
)

func init() {
	flag.Parse()
	log.Printf("fs: %s\naddr: %s\n", *flagFS, *flagAddr)

	var err error
	viewsFS, err = newFileSystem("views", "views")
	if err != nil {
		panic(err)
	}

	assetsFS, err = newFileSystem("assets", "assets")
	if err != nil {
		panic(err)
	}

	manager = newManager()
}

func newFileSystem(namespace, path string) (http.FileSystem, error) {
	switch *flagFS {
	case "packr":
		return packr.New(namespace, path), nil
	case "statik":
		return fs.NewWithNamespace(namespace)
	default:
		return http.Dir(path), nil
	}
}

func newManager() *views.Manager {
	// options
	opts := []views.Option{
		// views.Suffix(".tmpl"), // template suffix, default to .tmpl.
		// views.Delims("{{", "}}"), // template delimiters, default to "{{" and "}}".
		views.DefaultLayout("main"),
		views.LayoutsDir("layouts"),   // layout directory, relatived to views path.
		views.PartialsDir("partials"), // partials layout, relatived to layouts directory.
		// global function map for all templates.
		views.FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
		views.Cache(false), // disabled caching for developing.
	}
	manager = views.New(viewsFS, opts...)
	// add main layout.
	manager.AddLayout("main", "head", "header", "footer")
	// add a new layout.
	manager.AddLayout("page", "head")
	return manager
}

func main() {
	router := clevergo.NewRouter()
	router.Get("/", home)
	router.Get("/login", login)
	router.Get("/partial", partial)
	router.ServeFiles("/assets/*filepath", assetsFS)
	log.Println(http.ListenAndServe(*flagAddr, router))
}

func home(ctx *clevergo.Context) error {
	return manager.Render(ctx.Response, "site/index", map[string]interface{}{
		"title": "home",
	})
}

func login(ctx *clevergo.Context) error {
	return manager.RenderLayout(ctx.Response, "page", "user/login", nil)
}

func partial(ctx *clevergo.Context) error {
	return manager.RenderPartial(ctx.Response, "site/partial", map[string]interface{}{
		"title": "partial",
	})
}
