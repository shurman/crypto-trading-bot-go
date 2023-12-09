FROM        golang:alpine3.19 AS base
RUN         mkdir -p /app
WORKDIR     /app
COPY        . .
RUN         go mod download && go build -o ctbot

FROM        alpine:3.19
COPY        --from=base   /app/ctbot .
COPY        --from=base   /app/config.yml .
CMD         ["./ctbot"]