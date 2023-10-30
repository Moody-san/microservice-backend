pipeline {
    agent {
        docker {
            image 'moodysan/gobaseimage:latest' 
            args '--user root -v /var/run/docker.sock:/var/run/docker.sock' 
        }
    }
    environment {
        DOCKER_IMAGE = "moodysan/goapp:${BUILD_NUMBER}"
        GIT_REPO_NAME = "k8s-manifests"
    }
    stages {
        stage('Build go image') {
            steps {
                sh 'if [ -f go.mod ]; then rm -f go.mod; fi'
                sh 'go mod init example/app'
                sh 'go build -o build'
            }
        }
        stage('Build docker image and push to dockerHub') {
            steps {
                script {
                    sh 'docker build -t ${DOCKER_IMAGE} .'
                    withCredentials([usernamePassword(credentialsId: 'docker-cred', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                        sh 'docker login -u ${USERNAME} -p ${PASSWORD}'
                        sh 'docker push ${DOCKER_IMAGE}'
                    }
                }
            }
        }
        stage('Checkout Manifest Repo') {
            steps {
                git branch: 'main', url: 'https://github.com/Moody-san/k8s-manifests'
            }
        }
        stage('Update Manifest with newly create docker image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                    sh '''
                        ls
                        pwd
                        // cd && cd k8s-manifests/
                        // git config user.email "jenkins@gmail.com"
                        // git config user.name "jenkins"
                        // sed -i "s|\\(moodysan/goapp/\\).*|\\1${BUILD_NUMBER}|" deployment.yml
                        // git add deployment.yml
                        // git commit -m "Update deployment image to version ${BUILD_NUMBER}"
                        // git push https://${PASSWORD}@github.com/${USERNAME}/${GIT_REPO_NAME}.git HEAD:main
                    '''
                }
            }
        }
    }
}  
