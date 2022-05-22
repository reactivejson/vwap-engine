#!/usr/bin/env groovy

import com.nokia.neo.ci.constants.*

@Library('AANM_CI_tools_vars') _

def BuilderImage = "go-tools"
def GoVersion = "1.17.8"
def ToolsVersion = "67"

pipeline {
    agent {
        kubernetes {
            inheritFrom "${Neo.K8S_DIND_ROOTLESS}"
            yaml cciContainers([
                [
                    name: "go-container",
                    image: "${Docker.RELEASES_REPO}.${DOCKER_REGISTRY_URL}/${BuilderImage}:${GoVersion}-${ToolsVersion}",
                    command: ["cat"],
                    tty: "true"
                ],
                Neo.NEO_TOOLS_CONTAINER])
        }
    }
    environment {
        BUILD_VERSION = "0.0.${env.BUILD_NUMBER}"
        DOCKER_REGISTRY = "${Docker.INPROGRESS_REPO}.${DOCKER_REGISTRY_URL}"
        DOCKER_RELEASE_REGISTRY = "${Docker.RELEASES_REPO}.${env.DOCKER_REGISTRY_URL}"
        PACT_TAG = "${Neo.INPROGRESS_TAG}"
        IMAGE_NAME = "storage-engine"
        PROJECT_NAME = "storage-engine"
        APP_NAME = "storage-engine"
        COMMIT_ID = "${env.GIT_COMMIT}"
        DOCKER_HOST = sh(returnStdout: true, script: "echo tcp://\$(hostname -i):2375").trim()
    }
    options {
        timestamps()
        timeout(time: 15, unit: 'MINUTES')
    }

    triggers {
        gerrit(customUrl: '',
            gerritProjects: [[
                branches: [[
                    compareType: 'REG_EXP',
                    pattern: ".*"
                ]],
                compareType: 'PLAIN',
                disableStrictForbiddenFileVerification: false,
                pattern: "${GERRIT_PROJECT}",
            ]],
            serverName: "${GERRIT_NAME}",
            triggerOnEvents: [
                patchsetCreated(
                    excludeDrafts: true,
                    excludeNoCodeChange: false,
                    excludeTrivialRebase: false
                ),
                commentAddedContains('(?i)^(Patch Set [0-9]+:)?( [\\w\\\\+-]*)*(\\n\\n)?\\s*(recheck)')
            ]
        )
    }

    stages {
        stage("Setup environment") {
            parallel {
                stage('Start environment') {
                    steps {
                        container('go-container') {
                            sh 'make environment-start'
                        }
                    }
                }
                stage('Build') {
                    steps {
                        container('go-container') {
                            sh 'make build'
                        }
                    }
                }
                stage('Linting') {
                    steps {
                        container('go-container') {
                            sh 'make lint'
                        }
                    }
                }
            }
        }

        stage("Packaging and validating") {
            parallel {
                stage('Unit tests') {
                    steps {
                        container('go-container') {
                            sh 'make test'
                        }
                    }
                }
                stage('Integration tests') {
                    steps {
                        container('go-container') {
                            sh 'make test-integration'
                        }
                    }
                }
                stage('Create docker image and tags') {
                    steps {
                        container('go-container') {
                            sh 'make docker-build'
                        }
                    }
                }
                stage('Helm chart') {
                    steps {
                        container('go-container') {
                            sh 'make helm-chart'
                        }
                    }
                }
            }
        }
        stage('Sonar') {
            steps {
                script {
                    try {
                        sonarGoAnalysisWithReview(
                            // Build version
                            env.BUILD_VERSION,
                            // Project suffix
                            '',
                            // Sonar exclusions
                            '**/*_test.go,**/vendor/**,**/cmd/**',
                            // Sonar test exlusions
                            '**/vendor/**',
                            // Coverage report file
                            'target/test/*_coverage.txt',
                            // Golangci linter result file
                            'target/lint-report.xml'
                        )
                    } catch (Exception e) {
                        unstable('Sonar check failed')
                    }
                }
            }
        }
        stage('Publishing') {
            // All the containers are in the same workspace.
            // Be aware of thread safety when adding new parallel stages!
            parallel {
                stage('Publish docker image') {
                    steps {
                        container('go-container') {
                            withDockerRegistry([
                                credentialsId: "${Neo.ARTIFACTORY_CREDENTIALS_ID}",
                                url: "https://${Docker.INPROGRESS_REPO}.${DOCKER_REGISTRY_URL}"
                            ]) {
                                sh 'make docker-push'
                            }
                        }
                    }
                }
                stage('Publish helm chart') {
                    steps {
                        container('go-container') {
                            script {
                                buildInfo = Artifactory.newBuildInfo()
                                server = Artifactory.server env.ARTIFACTORY_SERVER_ID
                                server.credentialsId = Neo.ARTIFACTORY_CREDENTIALS_ID
                                def uploadSpec = """{
                                                        "files": [
                                                                    {
                                                                           "pattern": "builds/helm/*.tgz",
                                                                           "target": "${Helm.INPROGRESS_REPO}"
                                                                     }
                                                        ]
                                                    }"""
                                server.upload spec: uploadSpec, buildInfo: buildInfo, failNoOp: true
                            }
                        }
                    }
                }
            }
        }
        stage('Deploy') {
            environment {
                INVENTORY_PATH = ".inventory"
            }
            stages {
                stage('Resolve user lab and environment variables') {
                    steps {
                        script {
                            env.LAB = searchLab().toString()
                            echo 'Searched lab'
                            env.IMAGE_NAME = env.APP_NAME
                            env.CHART_NAME = env.IMAGE_NAME
                        }
                    }
                }
                stage('upgrade and rollback component') {
                    when {not {environment name: 'LAB', value: 'none'}}
                    steps {
                        upgradeComponent()
                    }
                    post {
                        always {
                            rollbackComponent()
                        }
                    }
                }
            }
        }
    }
    post {
        success {
            echo "Publishing build info to artifactory"
            script {
                buildInfo.env.collect()
                buildInfo.name = "${APP_NAME}"
                buildInfo.number = "${BUILD_VERSION}"
                server.publishBuildInfo buildInfo
            }
        }
        always {
            junit 'target/test/*-junit.xml'
        }
    }
}
