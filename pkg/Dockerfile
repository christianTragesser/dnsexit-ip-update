## container image ##
FROM gcr.io/distroless/static AS publish
COPY dist/dnsexit_linux_amd64_v1/dnsexit-linux-amd64 /usr/local/bin/dnsexit

ENTRYPOINT [ "dnsexit" ]