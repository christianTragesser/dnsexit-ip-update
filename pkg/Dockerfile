FROM golang:stretch AS source
WORKDIR $GOPATH/src/github.com/christiantragesser/dnsexit-ip-update
ADD go.mod .
ADD go.sum .
ADD main.go .
COPY dnsexit-ip-update ./dnsexit-ip-update

## linting env ##
FROM golangci/golangci-lint AS lint
WORKDIR $GOPATH/src/github.com/christiantragesser/dnsexit-ip-update
COPY --from=source /go/src/github.com/christiantragesser/dnsexit-ip-update .
RUN golangci-lint run

## unit test env ##
FROM source AS test
RUN go test -v --cover $GOPATH/src/github.com/christiantragesser/dnsexit-ip-update/dnsexit

## container image ##
FROM gcr.io/distroless/static AS publish
COPY dist/dnsexit_linux_amd64_v1/dnsexit-linux-amd64 /usr/local/bin/dnsexit

CMD [ "dispatch" ]