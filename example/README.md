Views Manager Example
---------------------

```shell
$ cd example
$ go run main.go

$ curl http://localhost:1234/
<html>
    <head>
        <title>Home</title>
    </head>
    <body>
        <header>Header</header>
        <h1>Hello World</h1>
        <footer>Footer</footer>
    </body>
</html

$ curl http://localhost:1234/login
<html>
    <head>
        <title>Login</title>
    </head>
    <body>
        <h1>Login</h1>
    </body>
</html>

$ curl http://localhost:1234/partial
<html>
    <head>
        <title>Partial</title>
    </head>
    <body>
        <h1>Partial</h1>
    </body>
</html>
```