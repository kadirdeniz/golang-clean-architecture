# AWS Deployment Guide

This guide explains how to set up and use the CI/CD pipeline for deploying the Go Clean Architecture application to AWS using GitHub Actions, ECR, and ECS.

## Table of Contents

- [Prerequisites](#prerequisites)
- [AWS Infrastructure Setup](#aws-infrastructure-setup)
- [GitHub Repository Configuration](#github-repository-configuration)
- [CI/CD Workflows](#cicd-workflows)
- [Deployment Process](#deployment-process)
- [Monitoring and Troubleshooting](#monitoring-and-troubleshooting)
- [Security Best Practices](#security-best-practices)

## Prerequisites

### Required Tools
- AWS CLI v2
- Docker
- Go 1.23+
- Git
- Make (optional, for using Makefile commands)

### AWS Account Setup
1. **AWS Account**: Ensure you have an AWS account with appropriate permissions
2. **IAM User**: Create an IAM user with programmatic access
3. **IAM Policies**: Attach the following managed policies to your IAM user:
   - `AmazonECS_FullAccess`
   - `AmazonEC2ContainerRegistryFullAccess`
   - `IAMFullAccess` (for creating roles)
   - `AmazonVPCFullAccess`
   - `ElasticLoadBalancingFullAccess`
   - `CloudWatchFullAccess`
   - `AmazonSSMFullAccess`

## AWS Infrastructure Setup

### 1. Create ECR Repository

```bash
# Create ECR repository for the application
aws ecr create-repository \
    --repository-name golang-clean-architecture \
    --region us-east-1

# Get the repository URI (save this for later)
aws ecr describe-repositories \
    --repository-names golang-clean-architecture \
    --region us-east-1 \
    --query 'repositories[0].repositoryUri' \
    --output text
```

### 2. Create ECS Infrastructure

You can use AWS CDK, CloudFormation, or Terraform to create the infrastructure. Here's a basic CloudFormation template:

```yaml
# infrastructure/ecs-infrastructure.yaml
AWSTemplateFormatVersion: '2010-09-09'
Description: 'ECS Infrastructure for Go Clean Architecture App'

Parameters:
  Environment:
    Type: String
    Default: dev
    AllowedValues: [dev, staging, prod]
  VpcId:
    Type: AWS::EC2::VPC::Id
    Description: VPC ID where resources will be created
  SubnetIds:
    Type: List<AWS::EC2::Subnet::Id>
    Description: Subnet IDs for the load balancer and ECS service

Resources:
  # ECS Cluster
  ECSCluster:
    Type: AWS::ECS::Cluster
    Properties:
      ClusterName: !Sub 'golang-clean-architecture-${Environment}-cluster'
      CapacityProviders:
        - FARGATE
        - FARGATE_SPOT
      DefaultCapacityProviderStrategy:
        - CapacityProvider: FARGATE
          Weight: 1

  # Application Load Balancer
  ApplicationLoadBalancer:
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
    Properties:
      Name: !Sub 'golang-clean-architecture-${Environment}-alb'
      Scheme: internet-facing
      Type: application
      Subnets: !Ref SubnetIds
      SecurityGroups:
        - !Ref LoadBalancerSecurityGroup

  # Security Groups
  LoadBalancerSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Security group for Application Load Balancer
      VpcId: !Ref VpcId
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 80
          ToPort: 80
          CidrIp: 0.0.0.0/0
        - IpProtocol: tcp
          FromPort: 443
          ToPort: 443
          CidrIp: 0.0.0.0/0

  ECSSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Security group for ECS tasks
      VpcId: !Ref VpcId
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 8080
          ToPort: 8080
          SourceSecurityGroupId: !Ref LoadBalancerSecurityGroup

  # Target Group
  TargetGroup:
    Type: AWS::ElasticLoadBalancingV2::TargetGroup
    Properties:
      Name: !Sub 'golang-clean-architecture-${Environment}-tg'
      Port: 8080
      Protocol: HTTP
      VpcId: !Ref VpcId
      TargetType: ip
      HealthCheckPath: /health
      HealthCheckProtocol: HTTP
      HealthCheckIntervalSeconds: 30
      HealthyThresholdCount: 2
      UnhealthyThresholdCount: 3

  # Listener
  LoadBalancerListener:
    Type: AWS::ElasticLoadBalancingV2::Listener
    Properties:
      DefaultActions:
        - Type: forward
          TargetGroupArn: !Ref TargetGroup
      LoadBalancerArn: !Ref ApplicationLoadBalancer
      Port: 80
      Protocol: HTTP

  # IAM Roles
  ECSTaskExecutionRole:
    Type: AWS::IAM::Role
    Properties:
      RoleName: !Sub 'ecsTaskExecutionRole-${Environment}'
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: ecs-tasks.amazonaws.com
            Action: sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
      Policies:
        - PolicyName: SSMParameterAccess
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action:
                  - ssm:GetParameter
                  - ssm:GetParameters
                  - ssm:GetParametersByPath
                Resource: !Sub 'arn:aws:ssm:${AWS::Region}:${AWS::AccountId}:parameter/${Environment}/*'

  ECSTaskRole:
    Type: AWS::IAM::Role
    Properties:
      RoleName: !Sub 'ecsTaskRole-${Environment}'
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: ecs-tasks.amazonaws.com
            Action: sts:AssumeRole

  # CloudWatch Log Group
  LogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      LogGroupName: !Sub '/ecs/golang-clean-architecture-${Environment}'
      RetentionInDays: 30

  # ECS Service
  ECSService:
    Type: AWS::ECS::Service
    DependsOn: LoadBalancerListener
    Properties:
      ServiceName: !Sub 'golang-clean-architecture-${Environment}-service'
      Cluster: !Ref ECSCluster
      LaunchType: FARGATE
      DeploymentConfiguration:
        MaximumPercent: 200
        MinimumHealthyPercent: 100
      DesiredCount: 1
      NetworkConfiguration:
        AwsvpcConfiguration:
          AssignPublicIp: ENABLED
          SecurityGroups:
            - !Ref ECSSecurityGroup
          Subnets: !Ref SubnetIds
      LoadBalancers:
        - ContainerName: golang-clean-architecture-app
          ContainerPort: 8080
          TargetGroupArn: !Ref TargetGroup
      TaskDefinition: !Ref TaskDefinition

  # Task Definition (initial placeholder)
  TaskDefinition:
    Type: AWS::ECS::TaskDefinition
    Properties:
      Family: !Sub 'golang-clean-architecture-task'
      NetworkMode: awsvpc
      RequiresCompatibilities:
        - FARGATE
      Cpu: 256
      Memory: 512
      ExecutionRoleArn: !Ref ECSTaskExecutionRole
      TaskRoleArn: !Ref ECSTaskRole
      ContainerDefinitions:
        - Name: golang-clean-architecture-app
          Image: !Sub '${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/golang-clean-architecture:latest'
          Essential: true
          PortMappings:
            - ContainerPort: 8080
              Protocol: tcp
          Environment:
            - Name: CONFIG_FILE_NAME
              Value: !Sub 'config.${Environment}'
            - Name: APP_ENV
              Value: !Ref Environment
          LogConfiguration:
            LogDriver: awslogs
            Options:
              awslogs-group: !Ref LogGroup
              awslogs-region: !Ref AWS::Region
              awslogs-stream-prefix: ecs

Outputs:
  ClusterName:
    Description: ECS Cluster Name
    Value: !Ref ECSCluster
    Export:
      Name: !Sub '${AWS::StackName}-ClusterName'
  
  ServiceName:
    Description: ECS Service Name
    Value: !Ref ECSService
    Export:
      Name: !Sub '${AWS::StackName}-ServiceName'
  
  LoadBalancerDNS:
    Description: Application Load Balancer DNS Name
    Value: !GetAtt ApplicationLoadBalancer.DNSName
    Export:
      Name: !Sub '${AWS::StackName}-LoadBalancerDNS'
```

Deploy the infrastructure:

```bash
# Deploy for development environment
aws cloudformation deploy \
    --template-file infrastructure/ecs-infrastructure.yaml \
    --stack-name golang-clean-architecture-dev \
    --parameter-overrides Environment=dev VpcId=vpc-xxxxxxxx SubnetIds=subnet-xxxxxxxx,subnet-yyyyyyyy \
    --capabilities CAPABILITY_NAMED_IAM \
    --region us-east-1

# Deploy for staging environment
aws cloudformation deploy \
    --template-file infrastructure/ecs-infrastructure.yaml \
    --stack-name golang-clean-architecture-staging \
    --parameter-overrides Environment=staging VpcId=vpc-xxxxxxxx SubnetIds=subnet-xxxxxxxx,subnet-yyyyyyyy \
    --capabilities CAPABILITY_NAMED_IAM \
    --region us-east-1

# Deploy for production environment
aws cloudformation deploy \
    --template-file infrastructure/ecs-infrastructure.yaml \
    --stack-name golang-clean-architecture-prod \
    --parameter-overrides Environment=prod VpcId=vpc-xxxxxxxx SubnetIds=subnet-xxxxxxxx,subnet-yyyyyyyy \
    --capabilities CAPABILITY_NAMED_IAM \
    --region us-east-1
```

### 3. Create SSM Parameters for Database Configuration

```bash
# Development environment
aws ssm put-parameter \
    --name "/dev/database/url" \
    --value "postgres://username:password@hostname:5432/database?sslmode=require" \
    --type "SecureString" \
    --region us-east-1

# Staging environment
aws ssm put-parameter \
    --name "/staging/database/url" \
    --value "postgres://username:password@hostname:5432/database?sslmode=require" \
    --type "SecureString" \
    --region us-east-1

# Production environment
aws ssm put-parameter \
    --name "/prod/database/url" \
    --value "postgres://username:password@hostname:5432/database?sslmode=require" \
    --type "SecureString" \
    --region us-east-1
```

## GitHub Repository Configuration

### 1. Repository Secrets

Add the following secrets to your GitHub repository (`Settings` > `Secrets and variables` > `Actions`):

- `AWS_ACCESS_KEY_ID`: Your AWS IAM user access key ID
- `AWS_SECRET_ACCESS_KEY`: Your AWS IAM user secret access key

### 2. Repository Variables (Optional)

You can add these as repository variables for easier management:

- `AWS_REGION`: `us-east-1` (or your preferred region)
- `ECR_REPOSITORY`: `golang-clean-architecture`

## CI/CD Workflows

The repository includes the following GitHub Actions workflows:

### 1. Continuous Integration (`ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**
- **Test and Code Quality**: Runs tests, linting, and code quality checks
- **Build and Push**: Builds Docker image and pushes to ECR (only for `main` and `develop` branches)

### 2. Development Deployment (`deploy-dev.yml`)

**Triggers:**
- Push to `develop` branch
- Manual dispatch

**Process:**
1. Downloads deployment artifacts from CI workflow
2. Deploys to development ECS service
3. Runs basic health checks

### 3. Staging Deployment (`deploy-staging.yml`)

**Triggers:**
- Push to `main` branch
- Manual dispatch

**Process:**
1. Requires manual approval
2. Downloads deployment artifacts
3. Runs pre-deployment smoke tests
4. Deploys to staging ECS service
5. Runs post-deployment tests

### 4. Production Deployment (`deploy-prod.yml`)

**Triggers:**
- Manual dispatch only

**Process:**
1. Requires strict input validation
2. Pre-deployment security checks
3. Manual approval with detailed information
4. Creates backup of current deployment
5. Deploys to production ECS service
6. Comprehensive health checks
7. Provides rollback instructions

### 5. Security Scanning (`security.yml`)

**Triggers:**
- Daily at 2 AM UTC
- Push to `main` or `develop` branches
- Pull requests
- Manual dispatch

**Scans:**
- Dependency vulnerability scanning
- Static code analysis
- Docker image security scanning
- License compliance checking
- Configuration security checks
- Secrets detection

### 6. Release Management (`release.yml`)

**Triggers:**
- Push tags matching `v*` pattern
- Manual dispatch

**Process:**
1. Validates release version
2. Builds for multiple architectures
3. Creates Docker image with version tags
4. Generates changelog
5. Creates GitHub release
6. Uploads binaries and checksums

## Deployment Process

### Development Deployment

1. **Automatic**: Push code to `develop` branch
2. **Manual**: Use GitHub Actions manual dispatch

```bash
# Push to develop branch (automatic deployment)
git checkout develop
git push origin develop
```

### Staging Deployment

1. **Automatic**: Merge `develop` to `main` branch
2. **Manual**: Use GitHub Actions manual dispatch

```bash
# Merge to main (triggers staging deployment)
git checkout main
git merge develop
git push origin main
```

The staging deployment requires manual approval. You'll receive a GitHub issue for approval.

### Production Deployment

Production deployments are **manual only** for safety:

1. Go to GitHub Actions
2. Select "Deploy to Production" workflow
3. Click "Run workflow"
4. Fill in the required inputs:
   - **Image tag**: The specific image tag to deploy (e.g., `abc123` or `v1.0.0`)
   - **Approved by**: Name of the person who approved the deployment
   - **Reason**: Reason for the production deployment

### Creating a Release

#### Automatic (Recommended)

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

#### Manual

1. Go to GitHub Actions
2. Select "Release Management" workflow
3. Click "Run workflow"
4. Fill in the version and release details

## Monitoring and Troubleshooting

### CloudWatch Logs

View application logs in CloudWatch:

```bash
# Get recent logs for development
aws logs tail /ecs/golang-clean-architecture-dev --follow

# Get recent logs for staging
aws logs tail /ecs/golang-clean-architecture-staging --follow

# Get recent logs for production
aws logs tail /ecs/golang-clean-architecture-prod --follow
```

### ECS Service Status

Check ECS service health:

```bash
# Development
aws ecs describe-services \
    --cluster golang-clean-architecture-dev-cluster \
    --services golang-clean-architecture-dev-service

# Staging
aws ecs describe-services \
    --cluster golang-clean-architecture-staging-cluster \
    --services golang-clean-architecture-staging-service

# Production
aws ecs describe-services \
    --cluster golang-clean-architecture-prod-cluster \
    --services golang-clean-architecture-prod-service
```

### Common Issues

#### 1. ECR Authentication Errors

```bash
# Login to ECR manually
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
```

#### 2. ECS Service Not Starting

- Check CloudWatch logs for error messages
- Verify task definition has correct image URI
- Ensure security groups allow traffic on port 8080
- Check IAM role permissions

#### 3. Load Balancer Health Check Failures

- Verify application responds on `/health` endpoint
- Check security group rules
- Ensure application binds to `0.0.0.0:8080`

### Rollback Procedures

#### Development/Staging Rollback

```bash
# Get previous task definition
aws ecs list-task-definitions --family-prefix golang-clean-architecture-task --status ACTIVE --sort DESC

# Update service to use previous task definition
aws ecs update-service \
    --cluster golang-clean-architecture-dev-cluster \
    --service golang-clean-architecture-dev-service \
    --task-definition golang-clean-architecture-task:PREVIOUS_REVISION
```

#### Production Rollback

The production deployment workflow provides specific rollback instructions after each deployment.

## Security Best Practices

### 1. GitHub Repository Security

- Enable branch protection rules for `main` and `develop` branches
- Require pull request reviews
- Enable required status checks
- Restrict force pushes
- Use environment protection rules

### 2. AWS Security

- Use least privilege IAM policies
- Enable CloudTrail for API logging
- Use AWS Secrets Manager or SSM Parameter Store for sensitive data
- Enable VPC Flow Logs
- Regularly rotate access keys

### 3. Application Security

- Use non-root user in Docker containers
- Scan images for vulnerabilities regularly
- Keep dependencies updated
- Use HTTPS/TLS for all communications
- Implement proper input validation

### 4. Secrets Management

- Never commit secrets to Git
- Use GitHub Secrets for CI/CD
- Use AWS SSM Parameter Store for application secrets
- Rotate secrets regularly
- Monitor secret access

## Environment-Specific Configurations

### Development
- **Purpose**: Testing new features and bug fixes
- **Deployment**: Automatic on push to `develop`
- **Resources**: Minimal (256 CPU, 512 MB RAM)
- **Database**: Shared development database

### Staging
- **Purpose**: Pre-production testing and validation
- **Deployment**: Manual approval required
- **Resources**: Production-like (512 CPU, 1024 MB RAM)
- **Database**: Staging database with production-like data

### Production
- **Purpose**: Live application serving real users
- **Deployment**: Strict manual process with approvals
- **Resources**: Production scale (1024+ CPU, 2048+ MB RAM)
- **Database**: Production database with backups and monitoring

## Cost Optimization

### Development
- Use FARGATE_SPOT for cost savings
- Lower CPU/memory allocations
- Single task instance

### Staging
- Balance cost and performance
- Scale down during off-hours
- Use smaller instance sizes

### Production
- Right-size based on actual usage
- Use Auto Scaling
- Monitor and optimize regularly
- Consider Reserved Instances for predictable workloads

## Maintenance and Updates

### Regular Tasks

1. **Weekly**: Review security scan results
2. **Monthly**: Update dependencies and base images
3. **Quarterly**: Review and optimize AWS costs
4. **As needed**: Update deployment configurations

### Dependency Updates

```bash
# Update Go dependencies
go get -u ./...
go mod tidy

# Update Docker base image in Dockerfile
# Check for Alpine/Go version updates

# Update GitHub Actions
# Review and update action versions in workflow files
```

### Infrastructure Updates

```bash
# Update CloudFormation stack
aws cloudformation deploy \
    --template-file infrastructure/ecs-infrastructure.yaml \
    --stack-name golang-clean-architecture-dev \
    --capabilities CAPABILITY_NAMED_IAM

# Apply changes to all environments
```

This comprehensive guide should help you successfully deploy and manage your Go Clean Architecture application on AWS using the provided CI/CD pipeline. 