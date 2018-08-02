#!/bin/bash
if [ -n "$1" ]
    then
        curl -u 'mapleque' https://api.github.com/user/repos -d '{"name":"'$1'"}'
        echo 'run command: git remote add origin git@github.com:mapleque/'$1'.git'
    else
        echo 'what your repos name?'
fi
