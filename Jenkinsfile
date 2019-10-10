#!/usr/bin/env groovy

node('swarm') {
    String APP_PATH = "/go/src/github.com/elastic/ecctl"

    stage('Checkout from GitHub') {
	    checkout scm
    }
    withCredentials([string(credentialsId: '2a9602aa-ab9f-4e52-baf3-b71ca88469c7', variable: 'GITHUB_TOKEN')]) {
        docker.image("golang:1.12-stretch").inside("-u root:root -v ${pwd()}:${APP_PATH} -w ${APP_PATH} -e APP_PATH=$APP_PATH -e GITHUB_TOKEN=$GITHUB_TOKEN") {
            try {
                stage("Download dependencies") {
                    // TODO: Remove the git config commands once it's been open sourced
                    sh 'git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"'
                    sh 'cd $APP_PATH && make vendor'
                    sh 'git config --global --unset url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf'
                }
                stage("Run Linters") {
                    sh 'cd $APP_PATH && make lint docs'
                }
                stage("Compile and Run Unit Tests") {
                    sh 'cd $APP_PATH && ./scripts/makewrap.sh unit unit-coverage'
                    publishHTML([allowMissing: false, alwaysLinkToLastBuild: false, keepAll: false, reportDir: 'reports', reportFiles: 'coverage.html', reportName: 'Coverage Report', reportTitles: ''])
                }
            } catch (Exception err) {
                throw err
            } finally {
                stage("Clean up") {
                    sh 'cd $APP_PATH && rm -rf reports dist bin /tmp/go-build-cache'
                }
            }
        }
    }
}
