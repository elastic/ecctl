#!/bin/bash

if [[ $(git status) == *"Your branch is behind"* ]]; then
    echo "=> Can't generate changelog: local branch is behind the remote master, please check out the latest changes."
    exit 1
fi

git fetch
PREV_TAG=$(git tag -l | tail -1)
CHANGELOGFILE=notes/${VERSION}.md

if [[ -z ${PREV_TAG} ]]; then echo "-> Exiting changelog generation since there's no previous tag"; exit 0; fi

echo "=> Attempting to generate the changelog for release ${VERSION}..."
read -p "=> Previous release was ${PREV_TAG}, is that correct? " -n 1 -r

if [[ ${REPLY} =~ ^[Yy]$ ]]; then
    echo ""
    cp scripts/changelog.tpl.md ${CHANGELOGFILE}
    git -c log.showSignature=false log --pretty="* %h %s" --no-decorate --no-color tags/${PREV_TAG}...master >> ${CHANGELOGFILE}
    
    echo "=> Changelog generated."
    VISUAL="${VISUAL:-"${EDITOR:-vim}"}"
    read -p "==> The changelog will be opened with ${VISUAL}, press enter to continue or specify your desired editor: " -r
    if [[ ${REPLY} != "" ]]; then
        VISUAL=${REPLY}
    fi
    "${VISUAL}" "${CHANGELOGFILE}"
fi


