def dir = "app3"
def deployments = [
    [branch: 'oracle', dirName: 'manifests-oracle', lbip: 'https://10.0.1.226:6443'],
    [branch: 'azure', dirName: 'manifests-azure', lbip: 'https://192.0.1.8:6443']
]
pipeline {
    agent any
    environment {
        DOCKER_ID = 'moodysan'
        DOCKER_PASSWORD = credentials('docker-password')
    }
    stages {
        stage('Checkout Application Repo') {
            steps {
                script {
                   echo 'clone application code from branch ${dir}....'
                   git branch: "${dir}", url: 'https://github.com/Moody-san/microservice-backend'
                }
            }
        }
        stage ('Build and Push Docker Image'){
            steps{
                script{
                    echo 'login to docker....'
                    sh "echo $DOCKER_PASSWORD | docker login -u $DOCKER_ID --password-stdin"
                    echo 'building image....'
                    sh "docker build -t ${DOCKER_ID}/${dir}:${BUILD_NUMBER} ."
                    sh "docker push ${DOCKER_ID}/${dir}:${BUILD_NUMBER}"
                }
            }
        }
        stage('Checkout and Update Manifest Repo') {
            steps {
                lock('deploymentlock'){
                    script {
                        deployments.each{ deployment ->
                            echo 'clone manifests repo....'
                            git branch: "${deployment.branch}", url: 'https://github.com/Moody-san/k8s-manifests'
                            echo "updating deployment file for ${deployment.dirName} cluster...."
                            withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                                def direxists = sh(script: "ls -1 ./manifests/ | grep ${dir}", returnStdout: true).trim()
                                if (!direxists.isEmpty()){
                                    sh """
                                        git config user.email "jenkins@mail.com"
                                        git config user.name "jenkins"
                                        sed -i "s|server: .*|server: ${deployment.lbip}|" appofapps/${dir}.yml
                                        git add appofapps/${dir}.yml
                                        sed -i "s|$DOCKER_ID/${dir}.*|${DOCKER_ID}/${dir}:${BUILD_NUMBER}|" manifests/${dir}/deployment.yml
                                        git add manifests/${dir}/deployment.yml
                                        git commit -m "Build number ${BUILD_NUMBER} deployed"
                                        git push https://${PASSWORD}@github.com/${USERNAME}/k8s-manifests.git HEAD:"${deployment.branch}"
                                    """
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
    options {
        skipDefaultCheckout()
        disableConcurrentBuilds()
    }
}  
