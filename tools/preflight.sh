#!/usr/bin/env sh

set -eu

cd "$(realpath "$(dirname "$(readlink -f "$0")")/..")"

if ! v=$(go version | grep -E -o "go1\.[^.]+"); then
    echo >&2 "Failed to identify go version"
    exit 1
fi

if test "$v" \< "go1.16"; then
    echo >&2 "Go version is too old: ${v}"
    echo >&2 "Go 1.16+ is required"
    exit 1
fi

if test \! \( -e /usr/include/sqlite3.h -o -e /usr/local/include/sqlite3.h \); then
    echo >&2 "Unable to find sqlite3.h header"
    echo >&2 "Please install libsqlite3 dev package"
    exit 1
fi

touch ".preflightcheck"
