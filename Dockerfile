FROM golang:1 as builder

COPY .  "/app/src"
WORKDIR "/app/src"

RUN make


FROM alpine:3

RUN apk add --no-cache ca-certificates tini

COPY --from=builder "/app/src/views/"                    "/views/"
COPY --from=builder "/app/src/honk"                      "/usr/local/bin/honk"
COPY --from=builder "/app/src/entrypoint.sh"             "/usr/local/bin/entrypoint"

WORKDIR "/var/empty"

ENV USER=""
ENV PASS=""
ENV ADDR=""

ENV PUID=1000
ENV PGID=1000

ENTRYPOINT ["tini", "-v", "--", "entrypoint"]
