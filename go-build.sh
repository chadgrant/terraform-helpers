#!/bin/bash
set -e

OUT="/out/dist"
BUILD_DATE=$(date -u '+%Y/%m/%d %H:%M:%S')

mkdir -p $OUT

if [ "$BUILD_NUMBER" != "" ]; then
  cd apply
  gox -ldflags="-X \"main.BuildDate=$BUILD_DATE\" -X \"main.Version=$BUILD_NUMBER\"" -output "$OUT/apply_{{.OS}}_{{.Arch}}"

  cd ../plan
  gox -ldflags="-X \"main.BuildDate=$BUILD_DATE\" -X \"main.Version=$BUILD_NUMBER\"" -output "$OUT/plan_{{.OS}}_{{.Arch}}"

  cd ../state
  gox -ldflags="-X \"main.BuildDate=$BUILD_DATE\" -X \"main.Version=$BUILD_NUMBER\"" -output "$OUT/state_{{.OS}}_{{.Arch}}"

  cd ../crypt
  gox -ldflags="-X \"main.BuildDate=$BUILD_DATE\" -X \"main.Version=$BUILD_NUMBER\"" -output "$OUT/crypt_{{.OS}}_{{.Arch}}"
fi

if [ "$GITHUB_TOKEN" != "" ]; then
  ghr -t $GITHUB_TOKEN -u $PROJECT_USERNAME -r $PROJECT_REPONAME -delete `git describe --tags` $OUT/
fi

chmod -R a+rw $OUT/.
