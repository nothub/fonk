#!/usr/bin/env sh

set -e

if ! git remote show | grep "upstream" > /dev/null; then
    # git hg bridge: https://github.com/felipec/git-remote-hg
    git remote add "upstream" "hg::https://humungus.tedunangst.com/r/honk"
fi

# merge upstream changes
git checkout master
git pull --log upstream master
