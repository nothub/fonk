FROM golang:1.20-bookworm as builder

RUN apt-get update -qy                           \
 && apt-get install -qy --no-install-recommends  \
      libsqlite3-dev                             \
 && apt-get clean -qy                            \
 && apt-get autoremove -qy                       \
 && rm -rf /var/lib/apt/lists/*

COPY .  "/app/src"
WORKDIR "/app/src"

RUN make


FROM debian:bookworm

RUN apt-get update -qy                           \
 && apt-get install -qy --no-install-recommends  \
    libsqlite3-dev                               \
    tini                                         \
 && apt-get clean -qy                            \
 && apt-get autoremove -qy                       \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder "/app/src/views/" "/usr/local/share/honk/views/"
COPY --from=builder "/app/src/honk"   "/usr/local/bin/honk"
COPY                "entrypoint.sh"   "/entrypoint.sh"

WORKDIR "/var/empty"

ENV USER=""
ENV PASS=""
ENV DOMAIN=""

ENV PUID=1000
ENV PGID=1000

ENTRYPOINT ["tini", "-v", "--", "/entrypoint.sh"]
