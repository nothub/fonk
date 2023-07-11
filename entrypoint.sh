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
    if  test -z "${USER}"   ||
        test -z "${PASS}"   ||
        test -z "${DOMAIN}"; then
        echo >&2 "missing env var for db init!"
        exit 1
    fi
    printf "%s\n%s\n%s\n%s\n" "${USER}" "${PASS}" "0.0.0.0:8080" "${DOMAIN}" | honk -datadir "/data" init
fi

honk -datadir "/data" -viewdir "/views" "$@"
