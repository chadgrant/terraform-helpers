#!/bin/bash

set -e

application="terraform-helpers"
out=`pwd`/out

rm -rf $out
mkdir -p $out

docker build -t ${application}_build .

docker run -it --rm \
       -v $out:/out/ \
       -e "BUILD_NUMBER=$BUILD_NUMBER" \
       -e "PROJECT_USERNAME=$CIRCLE_PROJECT_USERNAME" \
       -e "PROJECT_REPONAME=$CIRCLE_PROJECT_REPONAME" \
       -e "GITHUB_TOKEN=$GITHUB_TOKEN" \
       ${application}_build
