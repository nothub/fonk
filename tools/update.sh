#!/usr/bin/env sh

set -eu

cd "$(dirname "$(realpath "$0")")/.."

if ! type git-remote-hg >/dev/null 2>&1; then
    echo >&2 "Put git-remote-hg in your path to bridge git and hg:"
    echo >&2 "https://github.com/felipec/git-remote-hg"
    exit 1
fi

if ! git remote show | grep "upstream" >/dev/null 2>&1; then
    git remote add "upstream" "hg::https://humungus.tedunangst.com/r/honk"
fi

# merge upstream changes
git checkout master
git pull --log --no-rebase upstream master
