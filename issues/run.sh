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

DAYS_AGO="${3:-30}"
FROM_DATE=$(date -v -${DAYS_AGO}d -u +"%Y-%m-%dT%H:%M:%SZ")

curl "https://api.github.com/repos/$OWNER/$REPO/issues?since=$FROM_DATE" 2> /dev/null | \
jq 'map(select(.pull_request | not))' > issues.txt

jq length issues.txt