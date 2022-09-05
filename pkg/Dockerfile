FROM golang:alpine AS source
WORKDIR $GOPATH/src/github.com/christiantragesser/dnsexit-ip-update
ADD go.mod .
ADD go.sum .
ADD main.go .
COPY dnsexit ./dnsexit
RUN go mod tidy

## build ##
FROM source as linux-build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/dnsexit-linux-amd64 .

## container image ##
FROM gcr.io/distroless/static AS publish
COPY --from=linux-build /go/bin/dnsexit-linux-amd64 /usr/local/bin/dnsexit

ENTRYPOINT [ "dnsexit" ]