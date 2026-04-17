#!/bin/bash

# Exit early.
# See: https://www.gnu.org/savannah-checkouts/gnu/bash/manual/bash.html#The-Set-Builtin.
set -e

# Source dot env file.
source .env

# Arguments.
version=""
platform=""
dry_run=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --version)
      version="${2//v/}"
      shift 2
      ;;
    --platform)
      platform="$2"
      shift 2
      ;;
    --dry-run)
      dry_run="$2"
      shift 2
      ;;
    *)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

echo "Build and push 👷"
echo

echo "Gotenberg version: $version"
echo "Target platform: $platform"

if [ "$dry_run" = "true" ]; then
  echo "🚧 Dry run"
fi

IFS='/' read -ra arch <<< "$platform"
tag="$DOCKER_REGISTRY/$DOCKER_REPOSITORY:$version-${arch[1]}"

echo "Will use tag: $tag"
echo

# Build image.
run_cmd() {
  local cmd="$1"

  if [ "$dry_run" = "true" ]; then
    echo "🚧 Dry run - would run the following command:"
    echo "$cmd"
    echo
  else
    echo "⚙️ Running command:"
    echo "$cmd"
    eval "$cmd"
    echo
  fi
}

cmd="docker buildx build \
    --build-arg GOTENBERG_VERSION=$version \
    --platform $platform \
    --load \
    -t $tag \
    -f $DOCKERFILE $DOCKER_BUILD_CONTEXT
"
run_cmd "$cmd"

echo "✅ Done!"
echo "tags=$tag" >> "$GITHUB_OUTPUT"
exit 0
