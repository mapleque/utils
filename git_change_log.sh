#!/bin/bash


git describe --abbrev=0 --tags > /dev/null 2>&1
if [[ $? == 0 ]]; then
    git log --oneline --no-merges `git describe --abbrev=0 --tags`..HEAD|cut -c 9-|sort
else
    git log --format="- %ad %s" --no-merges --date=short
fi
