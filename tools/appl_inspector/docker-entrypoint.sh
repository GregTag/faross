#/bin/bash

# Error handling
echoerr() { echo "$@" 1>&2; }

error_exit()
{
    echoerr "$1"
    exit 1
}

if [ -z "$1" ]; then 
    error_exit "Pass purl as a first parameter"
fi

mkdir -p /tmp/pkg
oss-download --download-directory /tmp/pkg $1 1>/dev/null 2>&1 || error_exit "Error occurred while downloading the package"
mkdir /pkg
mv -t /pkg /tmp/pkg/*.tar.*z /tmp/pkg/*.zip 1>/dev/null 2>&1

RESULT_FILE="data.json"
appinspector analyze -o $RESULT_FILE -f json -s /pkg > /dev/null || error_exit "Error occurred while analyzing"

OUTPUT=()
# ------------------ Unsafe operations ------------------

TAGS_FILE="tags.txt"  # scorer_unsafe.py uses tags.txt
jq '.metaData.uniqueTags.[]' $RESULT_FILE > $TAGS_FILE
score=$(python3 scorer_unsafe.py) || error_exit "Error occurred in scorer_unsafe.py"
rm $TAGS_FILE

name="Unsafe operations"
desc="Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"
OUPUT+=("{\"checkName\": \"$name\", \"score\": $score, \"risk\": \"Medium\", \"description\": \"$desc\"}")

# ------------------ File types ------------------

TYPES_FILE="file_types.txt"  # scorer_file_types.py uses file_types.txt
jq '.metaData.fileExtensions.[]' $RESULT_FILE > $TYPES_FILE
score=$(python3 scorer_file_types.py) || error_exit "Error occurred in scorer_file_types.py"
rm $TYPES_FILE

name="File types"
desc="Determines the presence of executable file extensions"
OUPUT+=("{\"checkName\": \"$name\", \"score\": $score, \"risk\": \"High\", \"description\": \"$desc\"}")

# ------------------ Clean ------------------

printf '%s\n' "${OUPUT[@]}"
rm $RESULT_FILE
