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

TEMP_FILE='temp.txt'

issues_count=0
page=1
while : ; do
    curl "https://api.github.com/repos/$OWNER/$REPO/issues?since=$FROM_DATE&page=$page" 2> /dev/null > $TEMP_FILE
    received=$(jq length $TEMP_FILE)
    real_issues=$(jq 'map(select(.pull_request | not)) | length' $TEMP_FILE)
    ((issues_count+=real_issues))
    ((page++))
    [[ $received > 0 ]] || break
done

rm $TEMP_FILE

echo $issues_count