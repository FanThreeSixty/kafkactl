@Library('fanthreesixty') _

pipeline {
    agent {
        docker {
		label 'docker'
		image 'nexus.fanthreesixty.com:8082/golang:v1.0'
		registryUrl 'http://nexus.fanthreesixty.com:8082'
		registryCredentialsId 'NEXUS_DOCKER_REGISTRY'
		args '-v $WORKSPACE:/go/src/github.com/sporting-innovations/kafkactl'
        }
    }

    options {
        timestamps()
    }

    stages {
        stage('Checkout') {
            steps {
                scmCheckout()
            }
        }
        stage('Pre-Build Linting') {
            steps {
                ansiColor('xterm') {
			sh "go tool vet ./"
			sh "golint ./"
                }
            }
        }
        stage('Build') {
            steps {
                ansiColor('xterm') {
			sh "go build -v -o kafkactl-${env.BUILD_ID}-linux github.com/sporting-innovations/kafkactl"
                }
            }
        }
        stage('Test') {
            steps {
                ansiColor('xterm') {
                    echo "HAH.  Easy mode."
                }
            }
        }

        stage('Publish and Archive') {
            steps {
                ansiColor('xterm') {
                    echo "Yea, still not happening"
                    archiveArtifacts artifacts: "kafkactl-${env.BUILD_ID}*"
                }
            }
        }
    }
    post {
        always {
            //notifySlack channel: '#cloudops'
            //logstashSend failBuild: false, maxLines: 0
            cleanWs()
        }
    }
}
