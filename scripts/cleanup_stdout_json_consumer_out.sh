#!/bin/bash

set -e

source ./scripts/util.sh

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <input_file>"
  exit 1
fi

input_file="$1"

if [ ! -f "$input_file" ]; then
  echo "Error: File '$input_file' not found!"
  exit 1
fi

output_file="${input_file%.txt}.json"
clean_output_file="${output_file%.json}_clean.json"

# Remove escape characters before quotes in JSON property names
# Using 'sed' to replace '\"' with '"'
sed -i 's/\\"/"/g' "$input_file"

# Run the jq command to create an array of JSON objects.
# The input file usually looks like:
# {...}
# {...}
# And we want this to look like:
# [{...},{...}]
# To be valid JSON.
jq -s '.' "$input_file" > "$output_file"

# Now we need to transform a msg string to json object for even better results.
# Recursively clean the "msg" property in each object using jq.
# Attempts to parse the msg content as JSON (fromjson?). If parsing fails, the original content is kept (// .).
jq 'walk(if type == "object" and has("msg") then .msg |= (fromjson? // .) else . end)' "$output_file" > "$clean_output_file"

# Removed unused $output_file
rm -rf $output_file

if [ $? -eq 0 ]; then
  echo "Successfully cleaned up '$input_file' to '$clean_output_file'"
else
  echo "Error: Failed to process input file"
  exit 1
fi
