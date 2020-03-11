Views Manager Example
---------------------

## Structure

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

## Try it

```shell
$ cd example
$ go run main.go
```