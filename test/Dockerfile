FROM golang:1.22.2-alpine AS build-testrunner

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go build -o testrunner.exe main.go

FROM ubuntu AS build-spim

RUN apt update
RUN apt install --yes make build-essential
RUN apt install --yes bison flex

WORKDIR /app

COPY spim spim

WORKDIR /app/spim/spim

RUN make
RUN make install

RUN apt install --yes gdb

WORKDIR /app/test

COPY --from=build-testrunner /app/testrunner.exe /app/test/runner.exe

CMD ["/app/test/runner.exe"]

