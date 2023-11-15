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
                        expression { currentBuild.number == 14 }
                    }
                    steps {
                        script {
                            dir("apps"){
                                git branch: 'main', url: 'https://github.com/Moody-san/microservice-backend'
                                def changeddirsOutput = sh(script: "ls -l | awk '/^d/ {print \$9}'",returnStdout: true)
                                changeddirs=changeddirsOutput.split('\n').collect { it.trim() }.findAll { it }
                            }
                        }
                    }
                }
                stage('Add changed dirs to list'){
                    when {
                        expression { currentBuild.number != 14 }
                    }
                    steps {
                        script {
                            dir("apps"){
                                sh "git fetch origin main"
                                def changeddirsOutput = sh(script: "git diff --name-only main origin/main |cut -d/ -f1|uniq",returnStdout: true)
                                changeddirs=changeddirsOutput.split('\n').collect { it.trim() }.findAll { it }
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
                                        def dir = "${it}".trim()
                                        sh "git pull origin main:main"
                                        changeddirs.each(){
                                            dir("${dir}") {
                                                def image_name = "moodysan/${dir}:${BUILD_NUMBER}"
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
                                            def dir = "${it}".trim()
                                            sh """
                                                git config user.email "jenkins@gmail.com"
                                                git config user.name "jenkins"
                                                sed -i "s|moodysan/${dir}.*|moodysan/${dir}:${BUILD_NUMBER}|" services/${dir}/deployment.yml
                                                git add services/${dir}/deployment.yml
                                                git commit -m "Update ${dir} deployment image to version ${BUILD_NUMBER}"
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
            }
        }
    }
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }
}  
