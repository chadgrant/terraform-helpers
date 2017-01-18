#!/bin/bash
set -e

OUT="/out/dist"
BUILD_DATE=$(date -u '+%Y/%m/%d %H:%M:%S')
OSARCH="darwin/amd64 linux/amd64 windows/amd64 freebsd/amd64 linux/arm linux/arm64"

mkdir -p $OUT

if [ "$BUILD_NUMBER" != "" ]; then
  for d in "apply" "plan" "state" "crypt" "tfvars"
  do
    cd $d
    echo "Building $d..."
    gox -osarch "$OSARCH" \
      -ldflags="-X \"main.BuildDate=$BUILD_DATE\" -X \"main.Version=$BUILD_NUMBER\"" \
      -output "$OUT/${d}_{{.OS}}_{{.Arch}}"
    cd -
  done
fi

if [ "$GITHUB_TOKEN" != "" ]; then
  ghr -t $GITHUB_TOKEN -u $PROJECT_USERNAME -r $PROJECT_REPONAME -delete `git describe --tags` $OUT/
fi

chmod -R a+rw $OUT/.
