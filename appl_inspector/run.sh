# Error handling
error_exit()
{
    echo "$1"
    exit 1
}

if [ -z "$1" ]; then 
    error_exit "Pass path too directory to analyze as a first parameter"
fi
if [ ! -d $1 ]; then
    error_exit "Source directory does not exist"
fi

RESULT_FILE="data.json"
appinspector analyze -o $RESULT_FILE -f json -s $1 > /dev/null

# ------------------ Unsafe_operations ------------------
TAGS_FILE="tags.txt"  # scorer.py uses tags.txt
jq '.metaData.uniqueTags.[]' $RESULT_FILE > $TAGS_FILE
score=$(python3 scorer.py)
rm $TAGS_FILE

desc="Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"
echo "{\"Unsafe_operations\": {\"score\": $score, \"risk\": \"medium\", \"desc\": \"$desc\"}}"

# ------------------ Clean ------------------

rm $RESULT_FILE
