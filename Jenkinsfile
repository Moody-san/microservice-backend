def changeddirs = []
def directories = []
def buildall = "false"
pipeline {
    agent {
        docker {
            image 'moodysan/gobaseimage:latest'
            args '--user root -v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    stages {
        stage('Clean Workspace') {
            when {
                expression { currentBuild.number == 1 || "${buildall}"=="true"}
            }
            steps {
                script {
                    sh "echo clean workspace"
                    sh "rm -rf ./*"
                }
            }
        }
        stage('Checkout Application Repo') {
            when {
                expression { currentBuild.number == 1 || "${buildall}"=="true"}
            }
            steps {
                script {
                    dir("apps"){
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
        stage('Add Changed Directories to List'){
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
        stage ('Build Docker Images'){
            steps{
                script{
                    dir("apps"){
                        if (directories.size()>0){
                            sh "echo update local app repo with changes from remote origin"
                            sh "git pull origin main:main"
                            sh "echo building dockerfile for directories ${directories}"
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
            }
        }
        stage('Checkout Manifest Repo'){
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
                    dir("manifests"){
                        sh "echo update deployment files in manifests repo"
                        withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                            directories.each(){
                                sh """
                                    git config user.email "jenkins@gmail.com"
                                    git config user.name "jenkins"
                                    sed -i "s|moodysan/${it}.*|moodysan/${it}:${BUILD_NUMBER}|" manifests/${it}/deployment.yml
                                    git add manifests/${it}/deployment.yml
                                    git commit -m "Update ${it} deployment image to version ${BUILD_NUMBER}"
                                    git push https://${PASSWORD}@github.com/${USERNAME}/k8s-manifests.git HEAD:main
                                """
                            }
                        }
                    }
                }
            }
        }
        stage ('Remove tmp folders'){
            steps{
                script{
                    sh "echo remove tmp and manifest files generated recursively in workspace"
                    sh "rm -rf \$(find . -type d -name '*tmp*')"
                    sh "rm -rf manifests*"
                }
            }
        }
    }
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }
}  
