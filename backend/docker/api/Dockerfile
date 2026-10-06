FROM golang:1.27.1 AS builder

WORKDIR /srv
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/api ./cmd/api
RUN CGO_ENABLED=0 go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

FROM ubuntu:24.04

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends tzdata curl ca-certificates && rm -rf /var/lib/apt/lists/*
RUN dpkg-reconfigure -f noninteractive tzdata

WORKDIR /srv
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY ./tmp/api /srv/bin/api
COPY ./database/migrations /srv/migrations
COPY ./docker/api/entrypoint.sh /usr/local/bin/memebrary-entrypoint
RUN chmod +x /usr/local/bin/memebrary-entrypoint /srv/bin/api

ENV DJANGOLANG_API_ROOT=/api

HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 CMD /bin/sh -c 'curl -fsS "http://localhost:${PORT:-7070}/api/healthz" >/dev/null'
ENTRYPOINT ["/usr/local/bin/memebrary-entrypoint"]
CMD ["serve"]
