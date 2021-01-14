FROM alpine:latest

ARG VERSION

COPY dist/ssmwrap_v${VERSION}_linux_amd64.tar.gz /tmp/
RUN cd /tmp && \
    tar xzf ssmwrap_v${VERSION}_linux_amd64.tar.gz && \
    install ssmwrap_v${VERSION}_linux_amd64/ssmwrap /usr/local/bin/ssmwrap && \
    rm -rf ssmwrap_v${VERSION}_linux_amd64*
