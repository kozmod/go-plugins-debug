FROM golang:1.19.1-alpine3.15 AS build
ENV CGO_ENABLED=1
ENV GO111MODULE="on"

RUN apk update && apk add --no-cache musl-dev gcc libc-dev bash

WORKDIR /test_plugins
ADD . /test_plugins
RUN go build -gcflags "all=-N -l" -o out/main .
RUN go build -o out/plug1.so -gcflags="all=-N -l" -buildmode=plugin ./plugin/one/plugin1.go
RUN go build -o out/plug2.so -gcflags="all=-N -l" -buildmode=plugin ./plugin/two/plugin2.go

RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM alpine:3.15
EXPOSE 8080 40000
# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat
WORKDIR /
COPY --from=build /test_plugins/out/ /
COPY --from=build /go/bin/dlv /
# Run delve
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/main", "--", "2"]