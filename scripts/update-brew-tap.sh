#!/usr/bin/env bash

set -e

# Parameters
VERSION=$(echo ${1}| sed 's/^v//')

# Globals
FORMULA_FILE=/tmp/homebrew-tap/Formula/ecctl.rb
GITHUB_USER=marclop

# Execution

if [[ -d /tmp/homebrew-tap ]]; then rm -rf /tmp/homebrew-tap; fi
hub clone https://github.com/elastic/homebrew-tap /tmp/homebrew-tap

DARWIN_CHECKSUM=$(grep darwin_amd64 dist/*checksums.txt | awk '{print $1}' | tr -d '\n')
LINUX_CHECKSUM=$(grep linux_amd64.tar.gz dist/*checksums.txt | awk '{print $1}' | tr -d '\n')
OLD_DARWIN_CHECKSUM=$(grep sha256 ${FORMULA_FILE}|head -1| awk '{print $2}' | tr -d '\n' | tr -d '"')
OLD_LINUX_CHECKSUM=$(grep sha256 ${FORMULA_FILE}|tail -1| awk '{print $2}' | tr -d '\n' | tr -d '"')
OLD_VERSION=$(grep 'version \"' ${FORMULA_FILE} | awk '{print $2}' | tr -d '"' | tr -d '\n')

sed "s/${OLD_VERSION}/${VERSION}/g" ${FORMULA_FILE} | sed "s/${OLD_DARWIN_CHECKSUM}/${DARWIN_CHECKSUM}/" | sed "s/${OLD_LINUX_CHECKSUM}/${LINUX_CHECKSUM}/" > /tmp/ecctl.rb
mv /tmp/ecctl.rb ${FORMULA_FILE}

cd /tmp/homebrew-tap
hub fork --no-remote
hub remote add fork https://github.com/${GITHUB_USER}/homebrew-tap
hub add Formula/ecctl.rb
hub checkout -b f/update-ecctl-formula-to-${VERSION}
hub commit -m "Update ecctl version to ${VERSION}"
hub push fork f/update-ecctl-formula-to-${VERSION}
hub pull-request -m "Update ecctl version to ${VERSION}" -m "Created through automation by update-brew-tap.sh"
