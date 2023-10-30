pipeline {
  agent {
    docker {
      image 'a go docker image'
      args '--user root -v /var/run/docker.sock:/var/run/docker.sock' 
    }
  }
  stages {
    stage('Checkout Repo') {
        steps {

        }
    }
    stage('Build go image') {
        steps {

        }
    }
    stage('Build docker image and push to dockerHub') {
        steps {

        }
    }
    stage('Update Manifest with newly create docker image') {
        steps {

        }
    }
  }
}  