def directoryToImageMap = [:]
pipeline {
    agent none
    stages {
        stage('Build docker image and push to dockerHub') {
            agent {
                docker {
                    image 'moodysan/gobaseimage:latest'
                    args '--user root -v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
            steps {
                script {  
                    def directories = sh(script: 'ls -1 -d */', returnStdout: true).split('\n')
                    for (def dir in directories) {
                        dir = dir.replaceAll('/$', '')
                        def nochanges = sh(script: "git status -s ${dir} | grep -q ${dir}",returnStatus: true)
                        if (!nochanges) {
                            dir(dir) {
                                def image_name = "moodysan/${dir}:${BUILD_NUMBER}"
                                sh "docker build -t ${DOCKER_IMAGE} ."
                                def dockerImage = docker.image(image_name)
                                docker.withRegistry('https://registry.hub.docker.com','docker-cred') {
                                    dockerImage.push()
                                }
                                directoryToImageMap[dir] = image_name
                            }
                        } else {
                            sh "echo No changes detected in directory: ${dir}"
                        }
                    }
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
                        script {
                            withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                                for (dir in directoryToImageMap){
                                    sh '''
                                        git config user.email "jenkins@gmail.com"
                                        git config user.name "jenkins"
                                        sed -i "s|moodysan/${dir.key}.*|${dir.value}|" ${dir.key}/deployment.yml
                                        git add ${dir.key}/deployment.yml
                                        git commit -m "Update ${dir.key} deployment image to version ${BUILD_NUMBER}"
                                        git push https://${PASSWORD}@github.com/${USERNAME}/${GIT_REPO_NAME}.git HEAD:main
                                    '''
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}  
