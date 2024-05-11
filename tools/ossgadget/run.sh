# Error handling
echoerr() { echo "$@" 1>&2; }

error_exit()
{
    echoerr "$1"
    exit 1
}

if [ -z "$1" ]; then 
    error_exit "Pass purl to analyze as a first parameter"
fi

RESULT_FILE="data.json"
oss-detect-backdoor --format sarifv2 --output-file $RESULT_FILE $1 || error_exit "Error occurred while analyzing"

tags_matched=$(jq '.runs | first | .results | length' $RESULT_FILE)
score=$(( tags_matched < 10 ? 10 - tags_matched : 0 ))

name="Backdoors"
desc="Uses regular expressions to identify backdoors"
echo "{\"$name\": {\"score\": $score, \"risk\": \"medium\", \"desc\": \"$desc\"}}"

rm $RESULT_FILE