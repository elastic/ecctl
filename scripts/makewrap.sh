#!/usr/bin/env bash
set -e

export PROJECT="ecctl"
export FILE_DETAILS="$(ls -ln /go/src/github.com/elastic/${PROJECT}/Makefile)" ;
export JENKINS_UID=$(echo "${FILE_DETAILS}" | cut -d' ' -f3 -);
export JENKINS_GID=$(echo "${FILE_DETAILS}" | cut -d' ' -f4 -);


useradd -u "${JENKINS_UID}" elastic
chown -R ${JENKINS_UID}:${JENKINS_GID} /go/src/github.com/elastic/${PROJECT}
export GOCACHE='/tmp/go-build'
mkdir -p $(go env GOCACHE) && chmod -R 777 $(go env GOCACHE)
chown -R ${JENKINS_UID}:${JENKINS_GID} $(go env GOCACHE)
chown -R ${JENKINS_UID}:${JENKINS_GID} /go/pkg/

su - elastic -c "export GOCACHE='/tmp/go-build' && export PATH=\$PATH:~/usr/local/go/bin/ && make -C /go/src/github.com/elastic/${PROJECT} ${@}";
