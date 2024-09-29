# laxc-go

## User Guide
 ()
Make sure, that `go 1.22.0` or newer is properly installed on your system (https://go.dev/).

You can build the executable compiler using `go build -o laxc.exe main.go`.

You can run the compiler using `./laxc.exe <srcfile>` or `go run main.go <srcfile>`.

Explore the CLI using `./laxc.exe -h`.

You can look into the different compilation stages using `./laxc.exe <srcfile> -s <stage>`.
This will print out a textual representation of the specified stage output.

## Remarks

The compiler uses an intermediate language that is equivalent to but not compatible with the IL expected by the LaSCoT test platform.
In order to comply with tests, some transformations have to be done to the original IL.
You can view the original IL using `./laxc.exe <srcfile> -s intermediate`.

## Testing

Compile for test runner:

```
env GOOS=linux GOARCH=amd64 go build -o ../test/laxc.exe main.go
```