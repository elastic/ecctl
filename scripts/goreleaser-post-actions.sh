#!/usr/bin/env bash

set -e

# Parameters
VERSION=${1}

# Upload the binaries and the checksums.
for f in dist/*.{tar.gz,deb,rpm,txt}; do
    aws --profile ecsecurity s3 cp --acl bucket-owner-full-control \
    ${f} s3://download.elasticsearch.org/downloads/ecctl/$(echo ${VERSION}| sed 's/^v//')/
done

# Create the actual Github Release.
hub release create -F notes/${VERSION}.md -m ${VERSION}

# Update the brew tap formula
./scripts/update-brew-tap.sh ${VERSION}
