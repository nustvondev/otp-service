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
                    env.IMAGE_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT.take(6)}"
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
    }

    post {
        always {
            sh """
                echo "🧹 Cleaning up local image: ${FULL_IMAGE}"
                docker rmi ${FULL_IMAGE} || true
                docker image prune -f || true
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