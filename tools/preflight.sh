#!/usr/bin/env sh

set -eu

cd "$(realpath "$(dirname "$(readlink -f "$0")")/..")"

if ! v=$(go version | grep -E -o "go1\.[^.]+"); then
    echo >&2 "Failed to identify go version"
    exit 1
fi

if test "$v" \< "go1.20"; then
    echo >&2 "Go version is too old: ${v}"
    echo >&2 "Go 1.20+ is required"
    exit 1
fi

touch ".preflightcheck"
