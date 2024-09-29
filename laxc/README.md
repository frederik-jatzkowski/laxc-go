# laxc-go

Make sure, that `go 1.22.0` or newer is properly installed on your system.

You can build the executable compiler using:

```
go build -o laxc.exe main.go
```

Compile for test runner:

```
env GOOS=linux GOARCH=amd64 go build -o ../test/laxc.exe main.go
```

Get more help on how to use the compiler using

```
go run main.go --help
```