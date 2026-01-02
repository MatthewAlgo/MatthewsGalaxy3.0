// MatthewsGalaxy CI/CD Pipeline
// Comprehensive Jenkinsfile for enterprise deployments

pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '20'))
        timeout(time: 30, unit: 'MINUTES')
        timestamps()
        disableConcurrentBuilds()
    }

    environment {
        DOCKER_REGISTRY = 'ghcr.io'
        IMAGE_PREFIX = 'matthewsgalaxy'
        GO_VERSION = '1.21'
        NODE_VERSION = '20'
        
        // Credentials
        DOCKER_CREDENTIALS = credentials('docker-registry-credentials')
        SONAR_TOKEN = credentials('sonarqube-token')
        SLACK_WEBHOOK = credentials('slack-webhook-url')
    }

    parameters {
        choice(
            name: 'ENVIRONMENT',
            choices: ['staging', 'production'],
            description: 'Target deployment environment'
        )
        booleanParam(
            name: 'SKIP_TESTS',
            defaultValue: false,
            description: 'Skip test execution (not recommended)'
        )
        booleanParam(
            name: 'FORCE_DEPLOY',
            defaultValue: false,
            description: 'Force deployment even if no changes detected'
        )
    }

    stages {
        // ============================================
        // Stage: Checkout
        // ============================================
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    env.GIT_BRANCH = sh(
                        script: 'git rev-parse --abbrev-ref HEAD',
                        returnStdout: true
                    ).trim()
                    env.BUILD_VERSION = "${env.GIT_BRANCH}-${env.GIT_COMMIT_SHORT}-${env.BUILD_NUMBER}"
                }
                echo "Building version: ${env.BUILD_VERSION}"
            }
        }

        // ============================================
        // Stage: Install Dependencies
        // ============================================
        stage('Install Dependencies') {
            parallel {
                stage('Backend Dependencies') {
                    steps {
                        dir('backend') {
                            sh 'go mod download'
                            sh 'go mod verify'
                        }
                    }
                }
                stage('Email Service Dependencies') {
                    steps {
                        dir('email-service') {
                            sh 'go mod download'
                            sh 'go mod verify'
                        }
                    }
                }
                stage('Frontend Dependencies') {
                    steps {
                        dir('frontend') {
                            sh 'npm ci'
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Lint & Static Analysis
        // ============================================
        stage('Lint & Static Analysis') {
            when {
                expression { !params.SKIP_TESTS }
            }
            parallel {
                stage('Backend Lint') {
                    steps {
                        dir('backend') {
                            sh 'golangci-lint run --timeout=5m'
                        }
                    }
                }
                stage('Email Service Lint') {
                    steps {
                        dir('email-service') {
                            sh 'golangci-lint run --timeout=5m'
                        }
                    }
                }
                stage('Frontend Lint') {
                    steps {
                        dir('frontend') {
                            sh 'npm run lint'
                            sh 'npx tsc --noEmit'
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Unit Tests
        // ============================================
        stage('Unit Tests') {
            when {
                expression { !params.SKIP_TESTS }
            }
            parallel {
                stage('Backend Tests') {
                    steps {
                        dir('backend') {
                            sh 'go test -v -race -coverprofile=coverage.out -covermode=atomic ./...'
                            sh 'go tool cover -html=coverage.out -o coverage.html'
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'backend/coverage.html', allowEmptyArchive: true
                        }
                    }
                }
                stage('Email Service Tests') {
                    steps {
                        dir('email-service') {
                            sh 'go test -v -race -coverprofile=coverage.out -covermode=atomic ./...'
                        }
                    }
                }
                stage('Frontend Tests') {
                    steps {
                        dir('frontend') {
                            sh 'npm test -- --coverage --watchAll=false'
                        }
                    }
                    post {
                        always {
                            publishHTML(target: [
                                allowMissing: true,
                                alwaysLinkToLastBuild: true,
                                keepAll: true,
                                reportDir: 'frontend/coverage/lcov-report',
                                reportFiles: 'index.html',
                                reportName: 'Frontend Coverage Report'
                            ])
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Security Scanning
        // ============================================
        stage('Security Scan') {
            when {
                expression { !params.SKIP_TESTS }
            }
            parallel {
                stage('Go Security') {
                    steps {
                        sh 'gosec -no-fail -fmt json -out gosec-backend.json ./backend/...'
                        sh 'gosec -no-fail -fmt json -out gosec-email.json ./email-service/...'
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'gosec-*.json', allowEmptyArchive: true
                        }
                    }
                }
                stage('NPM Audit') {
                    steps {
                        dir('frontend') {
                            sh 'npm audit --audit-level=high || true'
                        }
                    }
                }
                stage('Container Scan') {
                    steps {
                        sh '''
                            docker build -t ${IMAGE_PREFIX}-backend:scan ./backend
                            trivy image --severity HIGH,CRITICAL --exit-code 0 ${IMAGE_PREFIX}-backend:scan
                        '''
                    }
                }
            }
        }

        // ============================================
        // Stage: Build Docker Images
        // ============================================
        stage('Build Docker Images') {
            parallel {
                stage('Build Backend') {
                    steps {
                        script {
                            docker.build("${DOCKER_REGISTRY}/${IMAGE_PREFIX}-backend:${BUILD_VERSION}", './backend')
                        }
                    }
                }
                stage('Build Email Service') {
                    steps {
                        script {
                            docker.build("${DOCKER_REGISTRY}/${IMAGE_PREFIX}-email-service:${BUILD_VERSION}", './email-service')
                        }
                    }
                }
                stage('Build Frontend') {
                    steps {
                        script {
                            docker.build("${DOCKER_REGISTRY}/${IMAGE_PREFIX}-frontend:${BUILD_VERSION}", './frontend')
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Push Docker Images
        // ============================================
        stage('Push Docker Images') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                    expression { params.FORCE_DEPLOY }
                }
            }
            steps {
                script {
                    docker.withRegistry("https://${DOCKER_REGISTRY}", 'docker-registry-credentials') {
                        sh """
                            docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-backend:${BUILD_VERSION}
                            docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-email-service:${BUILD_VERSION}
                            docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-frontend:${BUILD_VERSION}
                        """
                        
                        // Tag as latest for branch
                        if (env.GIT_BRANCH == 'main') {
                            sh """
                                docker tag ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-backend:${BUILD_VERSION} ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-backend:latest
                                docker tag ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-email-service:${BUILD_VERSION} ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-email-service:latest
                                docker tag ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-frontend:${BUILD_VERSION} ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-frontend:latest
                                docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-backend:latest
                                docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-email-service:latest
                                docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}-frontend:latest
                            """
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Deploy to Environment
        // ============================================
        stage('Deploy') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                    expression { params.FORCE_DEPLOY }
                }
            }
            stages {
                stage('Deploy to Staging') {
                    when {
                        anyOf {
                            branch 'develop'
                            expression { params.ENVIRONMENT == 'staging' }
                        }
                    }
                    steps {
                        script {
                            deployToEnvironment('staging')
                        }
                    }
                }
                stage('Deploy to Production') {
                    when {
                        anyOf {
                            branch 'main'
                            expression { params.ENVIRONMENT == 'production' }
                        }
                    }
                    steps {
                        timeout(time: 10, unit: 'MINUTES') {
                            input message: 'Deploy to production?', ok: 'Deploy'
                        }
                        script {
                            deployToEnvironment('production')
                        }
                    }
                }
            }
        }

        // ============================================
        // Stage: Health Check
        // ============================================
        stage('Health Check') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            steps {
                script {
                    def targetUrl = params.ENVIRONMENT == 'production' 
                        ? 'https://matthewsgalaxy.com'
                        : 'https://staging.matthewsgalaxy.com'
                    
                    retry(3) {
                        sleep(time: 10, unit: 'SECONDS')
                        sh "curl -f ${targetUrl}/health || exit 1"
                    }
                }
            }
        }
    }

    post {
        success {
            script {
                slackNotification('SUCCESS', 'green')
            }
            echo '✅ Pipeline completed successfully!'
        }
        failure {
            script {
                slackNotification('FAILURE', 'red')
            }
            echo '❌ Pipeline failed!'
        }
        unstable {
            script {
                slackNotification('UNSTABLE', 'yellow')
            }
        }
        cleanup {
            cleanWs()
            sh 'docker system prune -f || true'
        }
    }
}

// ============================================
// Helper Functions
// ============================================

def deployToEnvironment(String environment) {
    echo "Deploying to ${environment}..."
    
    def serverCredentials = environment == 'production' 
        ? 'production-server-ssh'
        : 'staging-server-ssh'
    
    def serverHost = environment == 'production'
        ? env.PRODUCTION_HOST
        : env.STAGING_HOST
    
    sshagent(credentials: [serverCredentials]) {
        sh """
            ssh -o StrictHostKeyChecking=no deploy@${serverHost} << 'EOF'
                cd /opt/matthewsgalaxy
                
                # Backup database before production deploy
                if [ "${environment}" = "production" ]; then
                    docker compose exec -T postgres pg_dump -U matthew matthewsgalaxy > backup_\$(date +%Y%m%d_%H%M%S).sql
                fi
                
                # Pull and deploy
                docker compose pull
                docker compose up -d --remove-orphans
                
                # Cleanup
                docker image prune -f
                
                echo "Deployment complete!"
EOF
        """
    }
}

def slackNotification(String status, String color) {
    def message = """
        *Matthew's Galaxy Pipeline ${status}*
        • Job: ${env.JOB_NAME}
        • Build: #${env.BUILD_NUMBER}
        • Branch: ${env.GIT_BRANCH}
        • Commit: ${env.GIT_COMMIT_SHORT}
        • Duration: ${currentBuild.durationString}
        • <${env.BUILD_URL}|View Build>
    """
    
    try {
        slackSend(
            color: color,
            message: message,
            channel: '#deployments'
        )
    } catch (Exception e) {
        echo "Slack notification failed: ${e.message}"
    }
}
