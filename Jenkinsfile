pipeline {
    agent any

    environment {
        DOCKERHUB_USER = "vjetnodejs2309"
        IMAGE_NAME     = "${DOCKERHUB_USER}/otp-service"
        IMAGE_TAG      = "${env.BUILD_NUMBER}-${env.GIT_COMMIT[0..5]}"
        FULL_IMAGE     = "${IMAGE_NAME}:${IMAGE_TAG}"
    }

    // Chỉ trigger khi push vào nhánh uat/*
    // Multibranch Pipeline tự filter theo cấu hình job
    triggers {
        githubPush()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo "▶ Branch  : ${env.BRANCH_NAME}"
                echo "▶ Commit  : ${env.GIT_COMMIT}"
                echo "▶ Image   : ${FULL_IMAGE}"
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
                    // 'dockerhub-creds' là ID credentials tạo trong Jenkins
                    docker.withRegistry('https://index.docker.io/v1/', 'dockerhub-creds') {
                        docker.image("${FULL_IMAGE}").push()
                    }
                }
            }
        }
    }

    post {
        always {
            // ✅ Clean image khỏi host sau khi push xong
            // || true để không fail pipeline nếu image đã bị xóa
            sh """
                echo "🧹 Cleaning up local image: ${FULL_IMAGE}"
                docker rmi ${FULL_IMAGE} || true
                docker image prune -f || true
            """
            cleanWs()
        }
        success {
            echo "✅ Done! Image pushed: ${FULL_IMAGE}"
        }
        failure {
            echo "❌ Build failed on branch: ${env.BRANCH_NAME}"
            // Gợi ý: thêm Slack/email notification ở đây sau
        }
    }
}