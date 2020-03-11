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
$ go get github.com/clevergo/views/example@master

// http.Dir file system.
$ $GOPATH/bin/example -addr=:4040

// packr file system.
$ $GOPATH/bin/example -fs=packr -addr=:4040

// statik file system.
$ $GOPATH/bin/example -fs=statik -addr=:4040
```