# Error handling
echoerr() { echo "$@" 1>&2; }

error_exit()
{
    echoerr "$1"
    exit 1
}

if [ -z "$1" ]; then 
    error_exit "Pass path to directory to analyze as a first parameter"
fi
if [ ! -d $1 ]; then
    error_exit "Source directory does not exist"
fi

RESULT_FILE="data.json"
appinspector analyze -o $RESULT_FILE -f json -s $1 > /dev/null || error_exit "Error occurred while analyzing"

# ------------------ Unsafe_operations ------------------
TAGS_FILE="tags.txt"  # scorer.py uses tags.txt
jq '.metaData.uniqueTags.[]' $RESULT_FILE > $TAGS_FILE
score=$(python3 scorer.py)
rm $TAGS_FILE

name="Unsafe operations"
desc="Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"
echo "{\"$name\": {\"score\": $score, \"risk\": \"medium\", \"desc\": \"$desc\"}}"

# ------------------ Clean ------------------

rm $RESULT_FILE
