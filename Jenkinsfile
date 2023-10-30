pipeline {
    agent {
        docker {
        image 'golang:1.17' 
        args '--user root -v /var/run/docker.sock:/var/run/docker.sock' 
        }
    }
    environment {
        DOCKER_IMAGE = "moodysan/goapp:${BUILD_NUMBER}"
        REGISTRY_CREDENTIALS = credentials('docker-cred')
        GITHUB_TOKEN = credentials('github-token')
        GIT_USER_NAME = "moody-san"
        GIT_REPO_NAME = "k8s-manifests"
    }
    stages {
        stage('Checkout Application Repo') {
            steps {
                git branch: 'main', url: 'https://github.com/Moody-san/microservice-backend'
            }
        }
        stage('Build go image') {
            steps {
                sh 'cd microservice-backend/ && go build -o build'
            }
        }
        stage('Build docker image and push to dockerHub') {
            steps {
                script {
                    sh 'cd microservice-backend/ && docker build -t ${DOCKER_IMAGE} .'
                    def dockerImage = docker.image("${DOCKER_IMAGE}")
                    docker.withRegistry('https://registry.hub.docker.com',"docker-cred") {
                        dockerImage.push()
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
                sh '''
                    cd && cd k8s-manifests/
                    git config user.email "jenkins@gmail.com"
                    git config user.name "jenkins"
                    sed -i "s|\\(moodysan/goapp/\\).*|\\1${BUILD_NUMBER}|" deployment.yml
                    git add deployment.yml
                    git commit -m "Update deployment image to version ${BUILD_NUMBER}"
                    git push https://${GITHUB_TOKEN}@github.com/${GIT_USER_NAME}/${GIT_REPO_NAME}.git HEAD:main
                '''
            }
        }
    }
}  
