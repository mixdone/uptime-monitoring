pipeline {
    agent { 
        docker {
            image 'golang:1.24'
        }
    }	

    environment {
        GOCACHE = '/tmp/go-build-cache'
    }

    stages {
        stage(build) {
            steps {
                sh 'go build -o main cmd/main.go'
                }
        }
            stage('Test') {
                steps {
                    sh 'go test ./...'
                }
            }
        }
}
