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

su-exec "honk" honk upgrade
su-exec "honk" honk --datadir "/data" --viewdir "/views" "$@"
