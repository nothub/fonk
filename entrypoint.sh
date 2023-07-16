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

honk --datadir "/data" --viewdir "/views" "$@"
