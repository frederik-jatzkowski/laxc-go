# Automatic Testing

The following commands must be executed from the `/test` directory.

## Building the Test Runner

Build the test runner using:

```
docker build -t relaxo-testrunner .
```
## Writing Tests

Tests are written as `txtar`-Archives (https://pkg.go.dev/golang.org/x/tools/txtar). A test is build like this:

```
Before the source, there can be comment lines.

-- src.lx --
<LAX-SRC-CODE>
-- compiler.exitstatus --
<DESIRED-EXIT-STATUS>
-- program.stdout --
<REGEX-TO-VALIDATE-PROGRAM-OUTPUT>
```

Strings are whitespace-trimmed for comparison.

All `.txtar`-files in the `/test/suite` directory and its subdirectories will be recognized as tests. So feel free to structure them according to your needs.

## Running Tests

Copy the binary of the to-be-tested LAX-compiler to `/test/laxc.exe`. The binary must be compiled for Ubuntu and the `x86_64` architecture. The test suite will be run against this compiler.

Run the test suite using:

```
docker run --rm -v ${PWD}/suite:/app/test/suite -v ${PWD}/laxc.exe:/app/test/laxc.exe relaxo-testrunner
```

Running with more output:

```
docker run --rm -v ${PWD}/suite:/app/test/suite -v ${PWD}/laxc.exe:/app/test/laxc.exe relaxo-testrunner verbose
```

Running in git bash on Windows:

```
docker run --rm -v /${PWD}/suite:/app/test/suite -v /${PWD}/laxc.exe:/app/test/laxc.exe relaxo-testrunner
```