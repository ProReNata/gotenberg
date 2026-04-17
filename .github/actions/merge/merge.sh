#!/bin/bash

# Exit early.
# See: https://www.gnu.org/savannah-checkouts/gnu/bash/manual/bash.html#The-Set-Builtin.
set -e

# Source dot env file.
source .env

# Arguments.
tags=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --tags)
      tags="$2"
      shift 2
      ;;
    *)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

echo "Merge tag(s) 👷"
echo

echo "Tag(s) to merge:"
IFS=',' read -ra tags_to_merge <<< "$tags"
for tag in "${tags_to_merge[@]}"; do
  echo "- $tag"
done
echo

# Build merge map.
declare -A merge_map

for tag in "${tags_to_merge[@]}"; do
  target_tag="${tag//-amd64/}"
  target_tag="${target_tag//-arm64/}"

  merge_map["$target_tag"]+="$tag "
done

# Merge tags.
for target in "${!merge_map[@]}"; do
  IFS=' ' read -ra source_tags <<< "${merge_map[$target]}"

  cmd="docker buildx imagetools create \
       -t $target \
       ${source_tags[*]}"

  echo "⚙️ Running command:"
  echo "$cmd"
  eval "$cmd"

  echo "➡️ $target pushed"
  echo
done

echo "✅ Done!"
exit 0
