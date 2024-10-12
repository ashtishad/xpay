# Docker Compose Guide for XPay

## Services Overview

1. postgres: PostgreSQL database service
2. api: XPay application service

## 1. postgres Service

```yaml
services:
  postgres:
    image: postgres:17.0-alpine3.20
    container_name: xpay_postgres
    environment:
      POSTGRES_DB: xpay
      POSTGRES_USER: ash
      POSTGRES_PASSWORD: samplepass
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ash -d xpay"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - xpay_network
```

### Purpose:
- Provides a PostgreSQL database for the XPay application.
- Uses a specific version (17.0) with Alpine Linux for a smaller footprint.
- Sets up initial database, user, and password.
- Exposes the standard PostgreSQL port.
- Persists data using a named volume.
- Implements a health check to ensure the database is ready before the api service starts.

### Best Practices:
1. Use specific image tags for consistency and reproducibility.
2. Implement health checks to ensure service readiness.
3. Use environment variables for configuration to keep sensitive data out of the compose file.
4. Use named volumes for data persistence.
5. Set a restart policy for improved reliability.

### Production/Staging Considerations:
- Use secret management for sensitive data (e.g., database passwords).
- Consider using a managed database service for production.
- Implement proper backup and recovery strategies.
- Adjust the exposed ports for better security (e.g., don't expose the database port publicly in production).

## 2. api Service

```yaml
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: xpay_app
    environment:
      DB_URL: "postgres://ash:samplepass@postgres:5432/xpay?sslmode=disable&timezone=UTC"
      SERVER_ADDRESS: "0.0.0.0:8080"
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - xpay_network
```

### Purpose:
- Builds and runs the XPay application.
- Configures the database connection and server address.
- Exposes the application port.
- Ensures the database is ready before starting.

### Best Practices:
1. Use build context and Dockerfile for application containerization.
2. Use environment variables for configuration.
3. Implement service dependencies with health checks.
4. Use a restart policy for improved reliability.


## Docker Compose in Production/Staging

While Docker Compose is primarily used for local development, the transition to AWS for staging and production environments involves a different set of tools and services. Here's a practical approach for both scenarios:

### Staging Environment on AWS

For staging, a simplified version of the production setup can be used:

1. **Infrastructure as Code**:
   - Use Terraform to define and manage AWS resources.
   - Example Terraform configuration for a basic staging setup:
     ```hcl
     resource "aws_ecs_cluster" "staging" {
       name = "xpay-staging-cluster"
     }

     resource "aws_ecs_task_definition" "xpay_api" {
       family                   = "xpay-api-staging"
       network_mode             = "awsvpc"
       requires_compatibilities = ["FARGATE"]
       cpu                      = "256"
       memory                   = "512"

       container_definitions = jsonencode([{
         name  = "xpay-api"
         image = "${aws_ecr_repository.xpay_api.repository_url}:latest"
         portMappings = [{
           containerPort = 8080
           hostPort      = 8080
         }]
       }])
     }

     resource "aws_ecr_repository" "xpay_api" {
       name = "xpay-api-staging"
     }

     resource "aws_rds_cluster" "postgresql" {
       cluster_identifier      = "xpay-staging-db"
       engine                  = "aurora-postgresql"
       engine_version          = "13.7"
       availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
       database_name           = "xpay"
       master_username         = "xpay_admin"
       master_password         = var.db_password
       backup_retention_period = 5
       preferred_backup_window = "07:00-09:00"
     }
     ```

2. **Continuous Deployment**:
   - Use AWS CodePipeline with CodeBuild for CI/CD.
   - Example buildspec.yml for CodeBuild:
     ```yaml
     version: 0.2

     phases:
       pre_build:
         commands:
           - echo Logging in to Amazon ECR...
           - aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_ENDPOINT
       build:
         commands:
           - echo Build started on `date`
           - echo Building the Docker image...
           - docker build -t $IMAGE_REPO_NAME:$IMAGE_TAG .
           - docker tag $IMAGE_REPO_NAME:$IMAGE_TAG $ECR_ENDPOINT/$IMAGE_REPO_NAME:$IMAGE_TAG
       post_build:
         commands:
           - echo Build completed on `date`
           - echo Pushing the Docker image...
           - docker push $ECR_ENDPOINT/$IMAGE_REPO_NAME:$IMAGE_TAG
           - printf '[{"name":"xpay-api","imageUri":"%s"}]' $ECR_ENDPOINT/$IMAGE_REPO_NAME:$IMAGE_TAG > imagedefinitions.json

     artifacts:
       files: imagedefinitions.json
     ```

3. **Database**:
   - Use Amazon RDS for PostgreSQL or Amazon Aurora PostgreSQL.
   - Implement a process for applying migrations as part of the deployment pipeline.

### Production Environment on AWS

For production, a more robust and scalable setup is required:

1. **ECS with Fargate**:
   - Deploy containers using ECS Fargate for serverless container management.
   - Use Application Load Balancer (ALB) for routing traffic.

2. **High Availability Database**:
   - Use Amazon Aurora PostgreSQL with read replicas for high availability and performance.
   - Implement automated backups and point-in-time recovery.

3. **Secrets Management**:
   - Use AWS Secrets Manager for storing and rotating sensitive information like database credentials.

4. **Monitoring and Logging**:
   - Set up CloudWatch for centralized logging and monitoring.
   - Use CloudWatch Alarms for critical metrics.
   - Consider implementing distributed tracing with AWS X-Ray.

5. **Security**:
   - Implement WAF (Web Application Firewall) rules on the ALB.
   - Use VPC with private subnets for ECS tasks and RDS instances.
   - Implement strict security groups and NACLs.

6. **Scaling**:
   - Set up Auto Scaling for ECS services based on CloudWatch metrics.
   - Use Aurora Auto Scaling for read replicas.

7. **Disaster Recovery**:
   - Implement multi-region deployments for critical components.
   - Set up regular disaster recovery drills.

Example Terraform configuration for a production ECS service:

```hcl
resource "aws_ecs_service" "xpay_api" {
  name            = "xpay-api-service"
  cluster         = aws_ecs_cluster.production.id
  task_definition = aws_ecs_task_definition.xpay_api.arn
  desired_count   = 3
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = aws_subnet.private[*].id
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.xpay_api.arn
    container_name   = "xpay-api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.front_end]
}

resource "aws_appautoscaling_target" "ecs_target" {
  max_capacity       = 10
  min_capacity       = 3
  resource_id        = "service/${aws_ecs_cluster.production.name}/${aws_ecs_service.xpay_api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "ecs_policy" {
  name               = "xpay-api-autoscaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.ecs_target.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs_target.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value = 70.0
  }
}
```
