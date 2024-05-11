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

OUTPUT=()
# ------------------ Unsafe operations ------------------

TAGS_FILE="tags.txt"  # scorer_unsafe.py uses tags.txt
jq '.metaData.uniqueTags.[]' $RESULT_FILE > $TAGS_FILE
score=$(python3 scorer_unsafe.py) || error_exit "Error occurred in scorer_unsafe.py"
rm $TAGS_FILE

name="Unsafe operations"
desc="Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"
OUPUT+=("{\"$name\": {\"score\": $score, \"risk\": \"medium\", \"desc\": \"$desc\"}}")

# ------------------ File types ------------------

TYPES_FILE="file_types.txt"  # scorer_file_types.py uses file_types.txt
jq '.metaData.fileExtensions.[]' $RESULT_FILE > $TYPES_FILE
score=$(python3 scorer_file_types.py) || error_exit "Error occurred in scorer_file_types.py"
rm $TYPES_FILE

name="File types"
desc="Determines the presence of executable file extensions"
OUPUT+=("{\"$name\": {\"score\": $score, \"risk\": \"high\", \"desc\": \"$desc\"}}")

# ------------------ Clean ------------------

printf '%s\n' "${OUPUT[@]}"
rm $RESULT_FILE
