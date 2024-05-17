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
CONFIDENT_MATCHES="matches.json"
TEMP_FILE="temp.json"
oss-detect-backdoor --format sarifv2 --output-file $RESULT_FILE $1 || error_exit "Error occurred while analyzing"
jq '.runs | first | .results | map( select(.properties.Confidence > 1) ) |
      map({id: .rule.id, severity: .properties.Severity, confidence: .properties.Confidence}) |
      unique' $RESULT_FILE > $CONFIDENT_MATCHES
jq 'map( 11 - .confidence - .severity)' $CONFIDENT_MATCHES > $TEMP_FILE

if [ $(jq '. | length' $CONFIDENT_MATCHES)  -eq 0 ]; then
    score=10
else
    start_score=$(jq '. | min' $TEMP_FILE)
    base_weigh=$(jq 'map( 8 - . ) | max' $TEMP_FILE)
    total_weight=$(jq 'map( 8 - . ) | add' $TEMP_FILE)
    subtr=$(( total_weight / base_weigh - 1))
    unbounded_score=$((start_score - subtr))
    score=$(( unbounded_score < 0 ? 0 : unbounded_score ))
fi

desc="Uses regular expressions to identify backdoors"
echo {\"checkName\": \"Backdoors\", \"score\": $score, \"risk\": \"Medium\", \"description\": \"$desc\"}

rm $RESULT_FILE
rm $TEMP_FILE
rm $CONFIDENT_MATCHES

# Table for start score
# confidence 4 (high)
# severity 2 => 5
# severity 4 => 3
# severity 6 => 1

# confidence 2 (medium)
# severity 2 => 7
# severity 4 => 5
# severity 6 => 3