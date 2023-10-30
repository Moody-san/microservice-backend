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
            sh 'cd microservice-backend/ && docker build -t ${DOCKER_IMAGE} .'
            def dockerImage = docker.image("${DOCKER_IMAGE}")
            docker.withRegistry('https://registry.hub.docker.com',"docker-cred") {
                dockerImage.push()
            }
        }
    }
    stage('Checkout Manifest Repo') {
        steps {
            git branch: 'main', url: 'https://github.com/Moody-san/k8s-manifest'
        }
    }
    stage('Update Manifest with newly create docker image') {
        steps {
            sh '''
                cd && cd k8s-manifest/
                git config user.email "jenkins@gmail.com"
                git config user.name "jenkins"
                BUILD_NUMBER=${BUILD_NUMBER}
                sed -i "s/replaceImageTag/${BUILD_NUMBER}/g" manifestrepository/deploy.yml
                git add manifestrepository/deploy.yml
                git commit -m "Update deployment image to version ${BUILD_NUMBER}"
                git push https://${GITHUB_TOKEN}@github.com/${GIT_USER_NAME}/${GIT_REPO_NAME} HEAD:main
            '''
        }
    }
  }
}  