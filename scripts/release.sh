#!/bin/bash
set -e

export BUILD_DATE=`date +%Y-%m-%d:%H:%M:%S`
export BUILD_VERSION=`git log -1 --pretty=%B | tr " " "\n" | grep -Ei 'v[0-9]+(\.[0-9]+)*'`

# export BUILDS="linux/amd64 darwin/amd64"


if [[ $BUILD_VERSION ]] 
then
  echo "Installing gox and ghr"
  go get github.com/mitchellh/gox
  go get github.com/tcnksm/ghr
  echo "Building $BUILD_VERSION"
  gox -ldflags "-X main.Version=${BUILD_VERSION} -X main.BuildDate=${BUILD_DATE}" -parallel=2 -output "dist/dfm_{{.OS}}_{{.Arch}}"

  ghr -t $GITHUB_TOKEN -u devctl -r devctl $BUILD_VERSION dist/

fi