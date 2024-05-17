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
oss-detect-backdoor --format sarifv2 --output-file $RESULT_FILE $1 2>/dev/null || error_exit "Error occurred while analyzing"
severity=$(jq '.runs | first | .results | map( select(.properties.Confidence > 1) ) |
     map({id: .rule.id, severity: .properties.Severity}) | unique |
     map(.severity) | add' $RESULT_FILE)
normalized_severity=$((severity / 2)) # деление с округлением вниз
score=$(( normalized_severity < 10 ? 10 - normalized_severity : 0 ))

desc="Uses regular expressions to identify backdoors"
echo {\"checkName\": \"Backdoors\", \"score\": $score, \"risk\": \"Medium\", \"description\": \"$desc\"}

rm $RESULT_FILE