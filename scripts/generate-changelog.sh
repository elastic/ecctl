#!/bin/bash

if [[ $(git status) == *"Your branch is behind"* ]]; then
    echo "=> Can't generate changelog: local branch is behind the remote master, please check out the latest changes."
    exit 1
fi

git fetch
PREV_TAG=$(git tag -l 'v[0-9-]*.[0-9-]*.[0-9-]'| tail -1)
CHANGELOGFILE=notes/${VERSION}.md
ADOC_CHANGELOG=docs/release_notes/${VERSION}.adoc
ADOC_CHANGELOG_HISTORY=docs/ecctl-release-notes.asciidoc

if [[ -z ${PREV_TAG} ]]; then echo "-> Exiting changelog generation since there's no previous tag"; exit 0; fi

MAJORMINOR="$(echo ${VERSION} | tr -d 'v' |cut -d '.' -f1).$(echo $VERSION | tr -d 'v' |cut -d '.' -f2)"

echo "=> Attempting to generate the changelog for release ${VERSION}..."
read -p "=> Previous release was ${PREV_TAG}, is that correct? " -n 1 -r

if [[ ${REPLY} =~ ^[Yy]$ ]]; then
    echo ""
    sed "s/VERSION_REPLACE/$(echo ${VERSION}| sed 's/^v//')/g" scripts/changelog.tpl.md | sed "s/VERSION_CROPPED_REPLACE/${MAJORMINOR}/" > ${CHANGELOGFILE}
    sed "s/VERSION_REPLACE/$(echo ${VERSION}| sed 's/^v//')/g" scripts/changelog.tpl.adoc > ${ADOC_CHANGELOG}
    
    git -c log.showSignature=false log --pretty="https://github.com/elastic/ecctl/commit/%h[%h] %s" --no-decorate --no-color tags/${PREV_TAG}..master |\
    sed 's|#\(.*\))|https://github.com/elastic/ecctl/pull/\1\[\#\1\])|' >> ${ADOC_CHANGELOG}
    date "+%n_Release date: %B %d, %Y_" >> ${ADOC_CHANGELOG}

    sed "5i\\
|* <<{p}-release-notes-${VERSION}>>" ${ADOC_CHANGELOG_HISTORY} | tr '|' '\n' > ${ADOC_CHANGELOG_HISTORY}.copy
    echo "include::release_notes/${VERSION}.adoc[]" >> ${ADOC_CHANGELOG_HISTORY}.copy
    mv ${ADOC_CHANGELOG_HISTORY}.copy ${ADOC_CHANGELOG_HISTORY}

    echo "=> Changelog generated."
    VISUAL="${VISUAL:-"${EDITOR:-vim}"}"
    read -p "==> The changelog (${ADOC_CHANGELOG}) will be opened with ${VISUAL}, press enter to continue or specify your desired editor: " -r
    if [[ ${REPLY} != "" ]]; then
        VISUAL=${REPLY}
    fi
    "${VISUAL}" "${ADOC_CHANGELOG}"
fi


