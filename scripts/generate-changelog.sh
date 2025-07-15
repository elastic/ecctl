#!/bin/bash

if [[ $(git status) == *"Your branch is behind"* ]]; then
    echo "=> Can't generate changelog: local branch is behind the remote master, please check out the latest changes."
    exit 1
fi

git fetch
PREV_TAG=$(git tag -l 'v[0-9-]*.[0-9-]*.[0-9-]' | sort --version-sort | tail -1)
CHANGELOGFILE=notes/${VERSION}.md

if [[ -z ${PREV_TAG} ]]; then echo "-> Exiting changelog generation since there's no previous tag"; exit 0; fi

MAJORMINOR="$(echo ${VERSION} | tr -d 'v' |cut -d '.' -f1).$(echo $VERSION | tr -d 'v' |cut -d '.' -f2)"

echo "=> Attempting to generate the changelog for release ${VERSION}..."
read -p "=> Previous release was ${PREV_TAG}, is that correct? " -n 1 -r

if [[ ${REPLY} =~ ^[Yy]$ ]]; then
    echo ""
    sed "s/VERSION_REPLACE/$(echo ${VERSION}| sed 's/^v//')/g" scripts/changelog.tpl.md | sed "s/VERSION_CROPPED_REPLACE/${MAJORMINOR}/" > ${CHANGELOGFILE}
fi


