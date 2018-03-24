#!/bin/sh

if [ -f "$(dirname "$0")/config.sh" ]; then
  . "$(dirname "$0")/config.sh"
fi

version=$1
DOCKER="${DOCKER:-sudo docker}"

usage() {
  echo 'Usage:'
  echo "  GITHUB_TOKEN=xxx $0 <version>"
  exit 1
}

if [ -z "${version}" ]; then
  usage
fi

if [ -z "${GITHUB_TOKEN}" ]; then
  echo 'GITHUB_TOKEN is empty'
  usage
fi

(set -x; $DOCKER run -it --rm \
  -v "$(pwd)":/usr/local/src/your-app skywinder/github-changelog-generator \
  --user lemonldap-ng-controller \
  --project lemonldap-ng-controller \
  --token "$GITHUB_TOKEN" \
  --future-release "$version"
)
