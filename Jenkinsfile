pipeline {
    agent none
    stages {
        stage('Build docker image and push to dockerHub') {
            agent {
                docker {
                    image 'moodysan/gobaseimage:latest'  
                }
            }
            environment {
                DOCKER_IMAGE = "moodysan/goapp:${BUILD_NUMBER}"
            }
            steps {
                script{
                    sh 'docker build -t ${DOCKER_IMAGE} .'
                    def dockerImage = docker.image("${DOCKER_IMAGE}")
                    docker.withRegistry('https://registry.hub.docker.com',"docker-cred") {
                        dockerImage.push()
                    }
                    sh 'ls -la'
                }
            }
        }
        stage('Run this outside docker container') {
            agent any
            environment {
                GIT_REPO_NAME = "k8s-manifests"
            }
            stages {
                stage('Get Manifest Repo'){
                    steps {
                        git branch: 'main', url: 'https://github.com/Moody-san/k8s-manifests'
                    }
                }
                stage('Update Manifest with newly create docker image') {
                    steps {
                        withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                            sh '''
                                git config user.email "jenkins@gmail.com"
                                git config user.name "jenkins"
                                sed -i "s|moodysan/goapp.*|moodysan/goapp:${BUILD_NUMBER}|" apps/goapp/deployment.yml
                                git add apps/goapp/deployment.yml
                                git commit -m "Update deployment image to version ${BUILD_NUMBER}"
                                git push https://${PASSWORD}@github.com/${USERNAME}/${GIT_REPO_NAME}.git HEAD:main
                            '''
                        }
                    }
                }
            }
            post {
                always {
                    cleanWs()
                }
            }
        }
    }
}  
