#!/bin/bash

# Error handling
error_exit()
{
    echo "$1"
    exit 1
}

if [ -z "$1" ]; then 
    error_exit "Pass OWNER as the first parameter"
fi
if [ -z "$2" ]; then 
    error_exit "Pass REPO as the second parameter"
fi

OWNER=$1
REPO=$2

curl "https://api.github.com/repos/$OWNER/$REPO" 2> /dev/null | jq ".stargazers_count"