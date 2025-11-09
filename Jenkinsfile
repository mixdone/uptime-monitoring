pipeline {
    agent { 
	docker {
	    image 'golang:1.24'
	}
    }	

    stages {
	stage(build) {
	    steps {
	    	sh 'go build -o maim cmd/main.go'
            }
	}
        stage('Test') {
            steps {
                sh 'go test'
            }
        }
    }
}
