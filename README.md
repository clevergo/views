# Views [![Build Status](https://travis-ci.org/clevergo/views.svg?branch=master)](https://travis-ci.org/clevergo/views) [![Coverage Status](https://coveralls.io/repos/github/clevergo/views/badge.svg?branch=master)](https://coveralls.io/github/clevergo/views?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/views)](https://goreportcard.com/report/github.com/clevergo/views) [![GoDoc](https://godoc.org/github.com/clevergo/views?status.svg)](http://godoc.org/github.com/clevergo/views)

Views is a templates(html/template) manager,  it provides the following features:

- **Nested template**: it is easy to create a complicated template, see [usage](#usage).
- **Cache**: allow to cache parsed templates(default to disable), see [benchmark](#benchmark).
- **Global settings**: it provides some useful setting for all templates, such as suffix, delimiters, funcMap etc. 

## Usage

```go
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
        // views.Cache(true),
    }
    view = views.New("./themes", opts...)
    // render template with layouts.
    view.Render(w, "site/index", data)
    // render tempalte without layouts.
    view.RenderPartial(w, "/site/partial")
```

Please take a look of the following [example](example):

### Example

```shell
$ cd example
$ go run main.go

$ curl http://localhost:1234/
<html>
    <head>
        <title>Home</title>
    </head>
    <body>
        <h1>Hello World</h1>
        <footer>I am footer</footer>
    </body>
</html>

$ curl http://localhost:1234/partial
<html>
    <head>
        <title>Standalone</title>
    </head>
    <body>
        <h1>Standalone</h1>
    </body>
</html>
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
