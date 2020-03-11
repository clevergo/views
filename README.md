# Views [![Build Status](https://travis-ci.org/clevergo/views.svg?branch=master)](https://travis-ci.org/clevergo/views) [![Coverage Status](https://coveralls.io/repos/github/clevergo/views/badge.svg?branch=master)](https://coveralls.io/github/clevergo/views?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/views)](https://goreportcard.com/report/github.com/clevergo/views) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/clevergo/views/v2) [![Release](https://img.shields.io/github/release/clevergo/views.svg?style=flat-square)](https://github.com/clevergo/views/releases)

Views is a templates(html/template) manager,  it provides the following features:

- **File System**: it use `http.FileSystem` to parse template files, allows to embed view files into go binary easilly by third-party tools, 
	such as [packr](https://github.com/gobuffalo/packr), [statik](https://github.com/rakyll/statik) etc. See [example](example).
- **Simple**: it bases on html/template, nothing more.
- **Cache**: allow to cache parsed templates(default to enabled), see [benchmark](#benchmark).
- **Global settings**: it provides some useful setting for all templates, such as suffix, delimiters, funcMap etc.

## Usage

```shell
$ go get github.com/clevergo/views/v2
```

Please take a look of the [example](example).

### Structure

Assume directory structure looks like:

```
views/
	layouts/                  contains layout files.
    	main.tmpl
    	page.tmpl
    	...
    	partials/             contains partial files.
        	head.tmpl
        	header.tmpl
        	footer.tmpl
        	...
	site/                     contains site's views.
	    home.tmpl
	    ...
	user/                     contains user's views.
    	login.tmpl
    	setting.tmpl
    	signup.tmpl
    	...
	...
```

### Initialize

```go
fs = http.Dir("./views") // file system, you can use packr, statik or other file system instead.
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
manager = views.New(viewsPath, opts...)
// add main layout.
manager.AddLayout("main", "head", "header", "footer")
// add a new layout.
manager.AddLayout("page", "head")
// add function to global funcMap
manager.AddFunc("foo", func() string {
    return "bar"
})
```

### Render

```go
// render with default layout
manager.Render(w, "site/index", nil)

// render with particular layout: page.
manager.RenderLayout(w, "page", "user/login", nil)

// render with data
manager.Render(w, "site/index", map[string]interface{}{
	"foo": "bar",
})

// render without layout.
manager.RenderPartial(w, "site/partial", nil)
```

## Benchmark

```shell
$ go test -bench=.
BenchmarkView_Render-12                     5670            235214 ns/op
BenchmarkView_RenderPartial-12             16450             64282 ns/op
BenchmarkCacheView_Render-12              237700              4724 ns/op
BenchmarkCacheView_RenderPartial-12       364700              3269 ns/op
```

The benchmark is base on the [example](example) that mentioned above, the result is depended on how complicated the template is. 
