def changeddirs = []
pipeline {
    agent none
    stages {
        stage('Following stages will run inside container') {
            agent {
                docker {
                    image 'moodysan/gobaseimage:latest'
                    args '--user root -v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
            stages{
                stage('Checkout Application Repo') {
                    when {
                        expression { currentBuild.number == 2 }
                    }
                    steps {
                        script {
                            dir("apps"){
                                git branch: 'main', url: 'https://github.com/Moody-san/microservice-backend'
                                changeddirs = sh(script: "ls -l | awk '/^d/ {print \$9}'|uniq",returnStdout: true).split('\n')
                            }
                        }
                    }
                }
                stage('Add changed dirs to list'){
                    when {
                        expression { currentBuild.number != 1 }
                    }
                    steps {
                        script {
                            dir("apps"){
                                sh "git fetch origin main"
                                changeddirs = sh(script: "git diff --name-only main origin/main |cut -d/ -f1|uniq",returnStdout: true).split('\n')
                            }
                        }
                    }
                }
                stage ('Checkout changes , build and push image'){
                    steps{
                        script{
                            try{
                                dir("apps"){
                                    if (!changeddirs.isEmpty()){
                                        sh "git pull origin main:main"
                                        changeddirs.each(){
                                            dir("${it}") {
                                                def image_name = "moodysan/${it}:${BUILD_NUMBER}"
                                                sh "docker build -t ${image_name} ."
                                                def dockerImage = docker.image("${image_name}")
                                                docker.withRegistry('https://registry.hub.docker.com','docker-cred') {
                                                    dockerImage.push()
                                                }
                                            }
                                        }
                                    }
                                    else{
                                        sh "echo No changes in any directories"
                                    }
                                }
                            }
                            catch (all){
                                sh "Unable to build docker image(s)"
                            }
                        }
                    }
                }
            }
        }
        stage('Following stages will run outside container') {
            agent any
            stages {
                stage('Get Manifest Repo'){
                    steps {
                        dir("manifests"){
                            git branch: 'main', url: 'https://github.com/Moody-san/k8s-manifests'
                        }
                    }
                }
                stage('Update Manifest with newly create docker image') {
                    steps {
                        script {
                            try{
                                dir("manifests"){
                                    withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                                        changeddirs.each(){
                                            sh """
                                                git config user.email "jenkins@gmail.com"
                                                git config user.name "jenkins"
                                                sed -i "s|moodysan/${it}.*|moodysan/${it}:${BUILD_NUMBER}|" services/${it}/deployment.yml
                                                git add $services/{it}/deployment.yml
                                                git commit -m "Update ${it} deployment image to version ${BUILD_NUMBER}"
                                                git push https://${PASSWORD}@github.com/${USERNAME}/k8s-manifests.git HEAD:main
                                            """
                                        }
                                    }
                                }
                            }
                            catch (all){
                                sh "Unable to update manifests"
                            }
                        }
                    }
                }
            }
        }
    }
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }
}  
