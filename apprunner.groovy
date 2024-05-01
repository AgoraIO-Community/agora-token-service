pipeline {
 agent { label 'ubuntu' }
  environment {
    GITHUB_TOKEN = credentials("CI_GITHUB_TOKEN")
    SERVICE_NAME = "wf-agora-token-server"
    APP_ID = credentials("wf-agora-app-id")
    APP_CERTIFICATE = credentials("wf-agora-app-cert")
  }
  parameters {
    string(name: 'SLACK_CHANNEL', defaultValue: 'bot-butlr-builds', description: 'Slack channel to notify')
  }
  options {
    timestamps ()
    disableConcurrentBuilds()
    }
  stages {
    stage('Clean') {
      steps {
        cleanWs()
      }
    }
    stage('Clone') {
        steps {
            sshagent(credentials: ['GITHUB_CI_SSH']){
                script {
                    checkout([
                    $class: 'GitSCM', branches: [[name: 'develop']], 
                    extensions: [
                        [$class: 'CleanBeforeCheckout', deleteUntrackedNestedRepositories: true], 
                        [$class: 'PruneStaleBranch']
                        ],
                    userRemoteConfigs: [[credentialsId: 'GITHUB_CI_SSH', url: 'git@github.com:WahooFitness/agora-token-service.git']]
                    ])
                }
            }
        }
    }
    stage('Deploy App') {
      steps {
        script {
          withCredentials([[
          $class: 'AmazonWebServicesCredentialsBinding',
          credentialsId: "jenkins_aws",
          accessKeyVariable: 'AWS_ACCESS_KEY_ID',
          secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
          ]]) {
          script {
            rc = sh(script: "chmod +x ./deploy-apprunner.sh && ./deploy-apprunner.sh", returnStatus: true)
            if(rc==0){
            env.APP_URL = sh(script: "aws apprunner list-services --query \"ServiceSummaryList[?contains(ServiceName, \'$SERVICE_NAME\')].ServiceUrl\" --output text --no-cli-pager", returnStdout: true).trim()
            }
          }
         }
        }
      }
    }
  }
  post {
    failure {
    script {
        slackSend channel: "${SLACK_CHANNEL}",
        color: 'bad',
        message: "The pipeline <${currentBuild.absoluteUrl}|${currentBuild.fullDisplayName}> failed on stage: ${STAGE_NAME}"
      }
    }
    cleanup {
        cleanWs()
    }
  }
}
