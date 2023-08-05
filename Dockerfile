FROM golang:alpine AS builder

WORKDIR /opt/apiportflux
COPY . .
RUN go build -o bin/main cmd/main.go

FROM alpine:latest

COPY --from=builder /opt/apiportflux/bin/main /opt/apiportflux/bin/main
COPY --from=builder /opt/apiportflux/config/*.json /opt/apiportflux/config/
COPY --from=builder /opt/apiportflux/infrastructure/postgres/migrations/*.sql /opt/apiportflux/infrastructure/postgres/migrations/

WORKDIR /opt/apiportflux

ARG version
ENV v $version

ENV env $environment
ENV p $port
ENV db $database
ENV dsn $dsn

CMD ["sh", "-c", "bin/main --ver $v --env $env --port $p --db $db --dsn $dsn"]