def dir = "app1"
def deployments = [
    [branch: 'oracle', dirName: 'manifests-oracle'],
    [branch: 'azure', dirName: 'manifests-azure'],
]
pipeline {
    agent {
        docker {
            image 'moodysan/gobaseimage:latest'
            args '--user root -v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    stages {
        stage('Checkout Application Repo') {
            steps {
                script {
                   sh "echo clone from branch ${dir} repository"
                   git branch: "${dir}", url: 'https://github.com/Moody-san/microservice-backend'
                }
            }
        }
        stage ('Build Docker Images'){
            steps{
                script{
                    dir("build"){
                        sh "echo building image"
                        def image_name = "moodysan/${dir}:${BUILD_NUMBER}"
                        docker build -t ${image_name} .
                        def dockerImage = docker.image("${image_name}")
                        docker.withRegistry('https://registry.hub.docker.com','docker-cred') {
                            dockerImage.push()
                        }
                    }
                }
            }
        }
        stage('Checkout and Update Manifest Repo') {
            steps {
                script {
                    deployments.each{ deployment ->
                        sh "echo updating deployment files for ${deployment.dirName} cluster"
                        dir("${deployment.dirName}"){
                            sh "echo clone manifests repo"
                            git branch: "${deployment.branch}", url: 'https://github.com/Moody-san/k8s-manifests'
                        }
                        dir("${deployment.dirName}"){
                            sh "echo update deployment files in manifests repo"
                            withCredentials([usernamePassword(credentialsId: 'GITHUB_TOKEN', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
                                directories.each(){
                                    sh """
                                        git config user.email "jenkins@mail.com"
                                        git config user.name "jenkins"
                                        sed -i "s|moodysan/${dir}.*|moodysan/${dir}:${BUILD_NUMBER}|" manifests/${dir}/deployment.yml
                                        git add manifests/${dir}/deployment.yml
                                        git commit -m "Update ${dir} deployment image to version ${BUILD_NUMBER}"
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
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }
}  
