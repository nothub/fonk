FROM golang:1 as builder

COPY .  "/app/src"
WORKDIR "/app/src"

RUN make


FROM alpine:3

RUN apk add --no-cache ca-certificates su-exec tini

COPY --from=builder "/app/src/entrypoint.sh"  "/usr/local/bin/entrypoint"
COPY --from=builder "/app/src/honk"           "/usr/local/bin/honk"
COPY --from=builder "/app/src/views/"         "/views/"

WORKDIR "/var/empty"

ENV PUID=1000
ENV PGID=1000

EXPOSE 8017/tcp
EXPOSE 8080/tcp

ENTRYPOINT ["tini", "-v", "--", "entrypoint"]
