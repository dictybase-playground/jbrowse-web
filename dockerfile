FROM golang:1.26.1-alpine3.23 AS go-builder

RUN apk update \
  && apk add git \
  && git clone https://github.com/caddyserver/caddy.git \
  && cd caddy/cmd/caddy \
  && go build -o /usr/local/bin/app 

FROM gcr.io/distroless/base-debian12 AS runner

COPY --from=go-builder /usr/local/bin/app /usr/local/bin/app
COPY ./jbrowse2 /www

WORKDIR /www
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/app", "file-server", "--listen", ":3000"]
