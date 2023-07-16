#!/usr/bin/env sh

set -e

addgroup                \
    --gid "${PGID}"     \
    --system            \
    honk

adduser                 \
    --uid "${PUID}"     \
    --system            \
    --ingroup honk      \
    --disabled-password \
    --no-create-home    \
    honk

mkdir -p "/data"
chown -R "honk:honk" "/data"

if test ! -f "/data/honk.db"; then
    echo >&2 "honk.db is missing, doing init..."
    if  test -z "${USER}" ||
        test -z "${PASS}" ||
        test -z "${ADDR}"; then
        echo >&2 "missing env var for db init!"
        exit 1
    fi
    honk --datadir "/data" init --username "${USER}" --password "${PASS}" --fqdn "${ADDR}"
    exit 0
fi

honk --datadir "/data" --viewdir "/views" "$@"
