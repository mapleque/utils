#!/bin/bash
git log --oneline --no-merges `git describe --abbrev=0 --tags`..HEAD|cut -c 9-|sort
