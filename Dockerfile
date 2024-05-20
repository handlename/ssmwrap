FROM alpine:latest

COPY dist/ssmwrap_linux_386/ssmwrap /usr/local/bin/ssmwrap

ENTRYPOINT ["/usr/local/bin/ssmwrap"]
