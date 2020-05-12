#!/usr/bin/env bash

set -e

# Parameters
VERSION=${1}
PROFILE=${2}

if [[ ! -z ${PROFILE} ]]; then
    PROFILE="--profile ${PROFILE}"
fi

# Upload the binaries and the checksums.
for f in dist/*.{tar.gz,deb,rpm,txt}; do
    aws ${PROFILE} s3 cp --acl bucket-owner-full-control \
    ${f} s3://download.elasticsearch.org/downloads/ecctl/$(echo ${VERSION}| sed 's/^v//')/
done

# Create the actual Github Release.
if [[ -z $(git config --get user.email) ]]; then
    git config user.email "cloud-delivery@elastic.co"
fi

if [[ -z $(git config --get user.name) ]]; then
    git config user.name "elasticcloudmachine"
fi

hub release create -F notes/${VERSION}.md ${VERSION}

# Update the brew tap formula
./scripts/update-brew-tap.sh ${VERSION}
