#!/bin/bash

SERVICE_NAME='wf-agora-token-server'
SERVICE_ARN=''
AWS_REGION='us-east-1'
ACCOUNT_ID='418541585715'
APPRUNNER_ACCESS_ROLE_NAME='AppRunnerECRAccessRole'
ECR_REPOSITORY="$SERVICE_NAME"
MAX_WAIT_TIME=900  # 15 minutes (900 seconds)
WAIT_INTERVAL=15   # 15 seconds

# Function to check and install Docker
check_docker() {
  if ! command -v docker &> /dev/null; then
    echo "Docker not found. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    sudo systemctl enable docker
    sudo systemctl start docker
    echo "Docker installed successfully."
  else
    echo "Docker is already installed."
  fi
}

# Function to check and install AWS CLI
check_awscli() {
  if ! command -v aws &> /dev/null; then
    echo "AWS CLI not found. Installing AWS CLI..."
    sudo apt-get update
    sudo apt-get install -y awscli
    echo "AWS CLI installed successfully."
  else
    echo "AWS CLI is already installed."
  fi
}

# Check if ECR repository exists
check_ecr_repository() {

# create instance configuration file
read -r -d '' ECR_LIFECYCLE_POLICY <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Expire images older than 14 days",
            "selection": {
                "tagStatus": "untagged",
                "countType": "sinceImagePushed",
                "countUnit": "days",
                "countNumber": 14
            },
            "action": {
                "type": "expire"
            }
        }
    ]
}
EOF
echo "${ECR_LIFECYCLE_POLICY}" > ecr_policy.json

    if aws ecr describe-repositories --repository-names $ECR_REPOSITORY >/dev/null 2>&1; then
      echo "ECR repository '$ECR_REPOSITORY' already exists."
    else
      echo "Creating ECR repository '$ECR_REPOSITORY'..."
      if ! aws ecr create-repository --repository-name $ECR_REPOSITORY --region $AWS_REGION --no-cli-pager; then
        echo "ECR creation failed"
        exit 1
      else
        echo "ECR created successfully"
        echo "Creating ECR Lifecycle policy rule"
        aws ecr put-lifecycle-policy --repository-name $ECR_REPOSITORY --lifecycle-policy-text file://ecr_policy.json
      fi
    fi
}

# Build container
build_push_container() {
  echo "Building and pushing Docker image..."
  if ! sudo docker build -t "$SERVICE_NAME:latest" --build-arg GITHUB_TOKEN="$GITHUB_TOKEN" --build-arg APP_ID=$APP_ID --build-arg APP_CERTIFICATE=$APP_CERTIFICATE --build-arg CORS_ALLOW_ORIGIN="*" .; then
    echo "Failed to build Docker image."
    exit 1
  fi

  if ! sudo docker tag "$SERVICE_NAME:latest" "$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$SERVICE_NAME:latest"; then
    echo "Failed to tag Docker image."
    exit 1
  fi

  if ! aws ecr get-login-password --region "$AWS_REGION" | sudo docker login --username AWS --password-stdin "$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"; then
    echo "Failed to log in to ECR."
    exit 1
  fi

  if ! sudo docker push "$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$SERVICE_NAME:latest"; then
    echo "Failed to push Docker image to ECR."
    exit 1
  fi

  echo "Docker image built and pushed to ECR."
}

# Create the App Runner service if it doesn't exist
create_app_runner_service_if_not_exists() {
  if check_app_runner_service; then
    update_app_runner_service
  else
    create_app_runner_service
  fi
}

# create serice configiration file
create_app_runner_service() {
  echo "Creating app runner service"
  if ! aws apprunner create-service \
      --service-name ${SERVICE_NAME} \
      --region ${AWS_REGION} \
      --cli-input-json file://apprunner_config.json \
      --no-cli-pager; then
    echo "App Runner service creation failed"
  else
    echo "App Runner service created successfully"
  fi
}

# update existing service
update_app_runner_service(){
  echo "Updating app runner service"
  if ! aws apprunner update-service \
      --service-arn ${SERVICE_ARN} \
      --region ${AWS_REGION} \
      --cli-input-json file://apprunner_config.json \
      --no-cli-pager; then
    echo "App Runner service update failed"
  else
    echo "App Runner service updating successfully"
  fi
}

# Check if App Runner service exists
check_app_runner_service() {
  service=$(aws apprunner list-services --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].{ARN:ServiceArn,Status:Status}" --output json --no-cli-pager)
  if [[ $service != "[]" ]]; then
    SERVICE_ARN=$(echo $service | jq -r '.[0].ARN')
    status=$(echo $service | jq -r '.[0].Status')
    echo "App Runner service '$SERVICE_NAME' exists."
    echo "Service ARN: $arn"
    echo "Service Status: $status"
    return 0
  elif [[ $status == "OPERATION_IN_PROGRESS" ]]; then
    wait_for_service_status
  else
    echo "App Runner service '$SERVICE_NAME' does not exist."
    return 1
  fi
}

# Wait for the service to reach the desired status
wait_for_service_status() {
  echo "Waiting for App Runner service '$SERVICE_NAME' to reach RUNNING status..."
  sleep 5
  start_time=$(date +%s)
  elapsed_time=0

  while [[ $elapsed_time -lt $MAX_WAIT_TIME ]]; do
    status=$(aws apprunner list-services --query "ServiceSummaryList[?contains(ServiceName, '$SERVICE_NAME')].Status" --output text --no-cli-pager)
    if [[ $status == "RUNNING" ]]; then
      echo "App Runner service '$SERVICE_NAME' is now RUNNING."
      export APP_URL=$(aws apprunner list-services --query "ServiceSummaryList[?contains(ServiceName, '$SERVICE_NAME')].ServiceUrl" --output text --no-cli-pager)
      return 0
    elif [[ "$status" == "CREATE_FAILED" ]]; then
            echo "Service deployment failed."
    else
      echo "Waiting for App Runner service '$SERVICE_NAME' to reach RUNNING status..."
    fi

    sleep $WAIT_INTERVAL

    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))
  done

  echo "Timeout: App Runner service '$SERVICE_NAME' did not reach the RUNNING state within $MAX_WAIT_TIME seconds."
  exit 1
}

function main() {

  check_docker
  check_awscli
  check_ecr_repository
  build_push_container
  create_app_runner_service_if_not_exists
  wait_for_service_status

}

main
