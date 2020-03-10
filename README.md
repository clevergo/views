# Views [![Build Status](https://travis-ci.org/clevergo/views.svg?branch=master)](https://travis-ci.org/clevergo/views) [![Coverage Status](https://coveralls.io/repos/github/clevergo/views/badge.svg?branch=master)](https://coveralls.io/github/clevergo/views?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/views)](https://goreportcard.com/report/github.com/clevergo/views) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/clevergo/views) [![Release](https://img.shields.io/github/release/clevergo/views.svg?style=flat-square)](https://github.com/clevergo/views/releases)

Views is a templates(html/template) manager,  it provides the following features:

- **Nested template**: it is easy to create a complicated template, see [usage](#usage).
- **Cache**: allow to cache parsed templates(default to enabled), see [benchmark](#benchmark).
- **Global settings**: it provides some useful setting for all templates, such as suffix, delimiters, funcMap etc.
- **[Hooks](#hooks)**: provides two hooks, `BeforeRenderEvent` and `AfterRenderEvent`.

## Usage

Please take a look of the [example](example).

### Structure

Let's take a look of an example views structure before digging into it:

```
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
manager.Render(w, "site/index", views.Context{
	"foo": "bar",
})

// render without layout.
manager.RenderPartial(w, "site/partial", nil)
```

## Hooks

```go
// regiters before render listener.
manager.RegisterOnBeforeRender(func(event *views.BeforeRenderEvent) error {
	return nil
})

// regiters after render listener.
manager.RegisterOnAfterRender(func(event *views.AfterRenderEvent) error {
	return nil
})
```

## Benchmark

```shell
$ go test -bench=.
BenchmarkView_Render-12                     6452            173335 ns/op
BenchmarkView_RenderPartial-12             31501             38441 ns/op
BenchmarkCacheView_Render-12              286412              4278 ns/op
BenchmarkCacheView_RenderPartial-12       316054              3750 ns/op
```

The benchmark is base on the [example](example) that mentioned above, the result is depended on how complicated the template is. 
