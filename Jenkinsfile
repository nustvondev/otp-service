pipeline {
    agent any

    environment {
        DOCKERHUB_USER = "vjetnodejs2309"
        IMAGE_NAME     = "${DOCKERHUB_USER}/otp-service"
    }

    triggers {
        githubPush()
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm

                script {
                    env.IMAGE_TAG  = "${env.BUILD_NUMBER}-${env.GIT_COMMIT.take(6)}"
                    env.FULL_IMAGE = "${env.IMAGE_NAME}:${env.IMAGE_TAG}"
                }

                echo "▶ Branch  : ${env.BRANCH_NAME}"
                echo "▶ Commit  : ${env.GIT_COMMIT}"
                echo "▶ Image   : ${env.FULL_IMAGE}"
            }
        }

        stage('Build Image') {
            steps {
                script {
                    docker.build("${FULL_IMAGE}", "--no-cache .")
                }
            }
        }

        stage('Push to DockerHub') {
            steps {
                script {
                    docker.withRegistry(
                        'https://index.docker.io/v1/',
                        'dockerhub-creds'
                    ) {
                        docker.image("${FULL_IMAGE}").push()
                    }
                }
            }
        }

        stage('Update ecom-platform') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'github-creds',
                    usernameVariable: 'GIT_USER',
                    passwordVariable: 'GIT_TOKEN'
                )]) {
                    sh """
                        git clone https://${GIT_USER}:${GIT_TOKEN}@github.com/nustvondev/ecom-platform.git ecom-platform

                        cd ecom-platform
                        sed -i 's/^OTP_SERVICE_TAG=.*/OTP_SERVICE_TAG=${IMAGE_TAG}/' .env
                        git config user.email "jenkins@ci.local"
                        git config user.name "Jenkins CI"
                        git add .env
                        git commit -m "ci: update otp-service tag to ${IMAGE_TAG}"
                        git push origin main
                    """
                }
            }
        }
    }

    post {
        always {
            sh """
                echo "🧹 Cleaning up local image: ${FULL_IMAGE}"
                docker rmi ${FULL_IMAGE} || true
                docker image prune -f || true
                rm -rf ecom-platform || true
            """

            deleteDir()
        }

        success {
            echo "✅ Done! Image pushed: ${FULL_IMAGE}"
        }

        failure {
            echo "❌ Build failed on branch: ${env.BRANCH_NAME}"
        }
    }
}