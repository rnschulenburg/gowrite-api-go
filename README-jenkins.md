pipeline {
agent any

environment {
APP_NAME = "gowrite-api-go"
INSTALL_PATH = "/usr/local/bin"
GO = "/usr/local/go/bin/go"
}

stages {

    stage('Checkout') {
      steps {
        git branch: 'main',
            url: 'https://github.com/rnschulenburg/gowrite-api-go.git'
      }
    }

    stage('Build') {
      steps {
        sh '''
        ${GO} mod tidy
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -o ${APP_NAME}
        '''
      }
    }

    stage('Deploy') {
      steps {
        sh '''
        sudo systemctl stop gowrite-api-go || true

        sudo mv ${APP_NAME} ${INSTALL_PATH}/${APP_NAME}
        sudo chown root:root ${INSTALL_PATH}/${APP_NAME}
        sudo chmod 755 ${INSTALL_PATH}/${APP_NAME}

        sudo systemctl daemon-reload
        sudo systemctl start gowrite-api-go
        sudo systemctl status gowrite-api-go --no-pager
        '''
      }
    }
}
}