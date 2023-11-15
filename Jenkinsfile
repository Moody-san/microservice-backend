def changeddirs = []
def directories = []
def cleanbuild = "true"
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
                stage('Clean & Checkout Application Repo') {
                    when {
                        expression { currentBuild.number == 1 || "${cleanbuild}"=="true"}
                    }
                    steps {
                        script {
                            dir("apps"){
                                sh "echo clean workspace"
                                sh "rm -rf ./*"
                                sh "echo cloning application repository"
                                git branch: 'main', url: 'https://github.com/Moody-san/microservice-backend'
                                sh "echo adding all directories to built"
                                changeddirs = sh(script: "ls -1 -l | awk '/^d/ {print \$9}'",returnStdout: true).split('\n')
                                for(def dir in changeddirs){
                                    if (!dir.contains('*tmp') && dir!=''){
                                        directories.add(dir)
                                    }
                                }
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
                                sh "echo update remote origin for application repo"
                                sh "git fetch origin main"
                                sh "echo adding directories that changed to list"
                                changeddirs = sh(script: "git diff --name-only main origin/main | grep '/' | cut -d/ -f1 | uniq",returnStdout: true).split('\n')
                                for(def dir in changeddirs){
                                    if (!dir.contains('*tmp') && dir!=''){
                                        directories.add(dir)
                                    }
                                }
                            }
                        }
                    }
                }
                stage ('Checkout changes , build and push image'){
                    steps{
                        script{
                            try{
                                dir("apps"){
                                    if (directories.size()>0){
                                        sh "echo update local app repo with changes from remote origin"
                                        sh "git pull origin main:main"
                                        sh "echo building dockerfile for directories -> ${directories}"
                                        directories.each(){
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
                            catch (Exception e) {
                                echo 'Exception occurred: ' + e.toString()
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
                            sh "echo clone manifests repo"
                            git branch: 'main', url: 'https://github.com/Moody-san/k8s-manifests'
                        }
                    }
                }
                stage('Update Manifest with newly create docker image') {
                    steps {
                        script {
                            try{
                                dir("manifests"){
                                    sh "echo update deployment files in manifest repo"
                                    withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                                        directories.each(){
                                            sh """
                                                git config user.email "jenkins@gmail.com"
                                                git config user.name "jenkins"
                                                sed -i "s|moodysan/${it}.*|moodysan/${it}:${BUILD_NUMBER}|" services/${it}/deployment.yml
                                                git add services/${it}/deployment.yml
                                                git commit -m "Update ${it} deployment image to version ${BUILD_NUMBER}"
                                                git push https://${PASSWORD}@github.com/${USERNAME}/k8s-manifests.git HEAD:main
                                            """
                                        }
                                    }
                                }
                            }
                            catch (Exception e) {
                                echo 'Exception occurred: ' + e.toString()
                            }
                        }
                    }
                }
                stage ('Remove tmp folders'){
                    steps{
                        script{
                            sh "echo remove tmp files generated recursively in and outside workspace"
                            sh "rm -rf \$(find . -type d -name '*tmp*')"
                            sh "rm -rf \$(find .. -type d -name '*tmp*')"
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
