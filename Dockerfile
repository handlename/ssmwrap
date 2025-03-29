FROM alpine:latest

COPY dist/ssmwrap_linux_386_sse2/ssmwrap /usr/local/bin/ssmwrap

ENTRYPOINT ["/usr/local/bin/ssmwrap"]
