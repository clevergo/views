Views Manager Example
---------------------

```shell
$ go get github.com/clevergo/views/example@master

// http.Dir file system.
$ $GOPATH/bin/example -addr=:4040

// packr file system.
$ $GOPATH/bin/example -fs=packr -addr=:4040

// statik file system.
$ $GOPATH/bin/example -fs=statik -addr=:4040
```

And then visit http://localhost:4040/.