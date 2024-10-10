# Configuration Management Guide

## General Setup

In all environments, we use a YAML configuration file named `config.yaml`. The application looks for this file in the project root directory.

### Configuration File Structure

Create a `config.yaml` file in the project root with the following structure:

```yaml
app:
  env: dev  # Options: dev, staging, production
  gin_mode: release # Options: debug, release
  server_address: ":8080"

db:
  url: "postgres://user:password@host:port/dbname?sslmode=disable&timezone=UTC"
  max_open_conns: 18
  max_idle_conns: 18
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"

jwt:
  private_key: "base64_encoded_private_key"
  public_key: "base64_encoded_public_key"

card:
  aes_key: "base64_encoded_aes_key"
```

## Local Development

### Setup:
1. Copy `local-dev/config.yaml.example` to `config.yaml` in the project root.
2. Adjust the settings in `config.yaml` for your local environment.
3. Run the application with `make run`.

### Best Practices:
- Keep `config.yaml` out of version control (add to `.gitignore`).
- Use dummy values for sensitive data in local development.

## Staging Environment

### Setup:
1. Create a staging-specific `config.yaml` file.
2. Use environment variables to override sensitive values if needed.
3. Deploy the `config.yaml` file to the project root in your staging environment.

### Best Practices:
- Use a staging-specific secrets manager for sensitive data.
- Use reduced-privilege credentials.
- Set up automated deployment pipelines to ensure consistency.

## Production Environment

### Setup:
1. Create a production `config.yaml` file with non-sensitive configurations.
2. Use environment variables or a secret management service for sensitive data.
3. For Kubernetes: Use ConfigMap to deploy `config.yaml` to the project root and Secrets for sensitive data.

### Best Practices:
- Use cloud provider's secret management service for sensitive data.
- Implement least-privilege access for services and users.
- Use immutable infrastructure and automated deployments.
- Regularly rotate secrets and access keys.

## General Best Practices

1. Configuration Validation:
   - Validate all required configurations on application startup.
   - Fail fast if any critical configuration is missing.

2. Logging:
   - Log non-sensitive configuration values for debugging.
   - Never log sensitive information.

3. Secrets Management:
   - Use a dedicated secrets manager for highly sensitive data.
   - Implement secret rotation and versioning.

4. Access Control:
   - Implement strict access controls for configuration management.
   - Use audit trails for configuration changes.

5. Documentation:
   - Maintain up-to-date documentation of all configuration parameters.
   - Include clear instructions for setup in different environments.

6. Environment-Specific Configurations:
   - Create separate `config.yaml` files for each environment (dev, staging, prod).
   - Ensure your deployment process places the correct `config.yaml` in the project root.

7. Version Control:
   - Store template or example `config.yaml` files in version control.
   - Use these templates to generate environment-specific configurations during deployment.

Remember: Never commit sensitive information to version control. Always use environment-specific configuration management in your CI/CD pipeline.

## Overriding Configurations

While `config.yaml` is the primary configuration source, you can override specific values using environment variables. This is particularly useful for sensitive data, environment-specific settings, or when running the application in containerized environments.

### Environment Variable Override

You can override configurations by setting environment variables directly. The application uses Viper's automatic environment variable binding. For nested YAML keys, use underscores to separate levels. For example:

```sh
export APP_SERVER_ADDRESS=0.0.0.0:8080
export DB_URL="new_connection_string"
```

These environment variables will take precedence over the values in `config.yaml`.

### Docker and Docker Compose Considerations

When running the application using Docker or Docker Compose, you may need to adjust certain configurations, particularly the database connection string. For example, the hostname for the PostgreSQL database typically changes from "localhost" to the service name defined in your Docker Compose file (often "postgres").

Here's an example of how you might override configurations when running the Docker container:

```makefile
docker-run:
	@docker run --name xpay_app --network xpay_network \
		-e DB_URL="postgres://ash:lol@postgres:5432/xpay?sslmode=disable&timezone=UTC" \
		-e APP_SERVER_ADDRESS="0.0.0.0:8080" \
		-e APP_GIN_MODE="release" \
		-p 8080:8080 xpay:latest
```

In this example:
- `DB_URL` is modified to use "postgres" as the hostname instead of "localhost".
- `APP_SERVER_ADDRESS` is set to "0.0.0.0:8080" to allow connections from outside the container.
- `APP_GIN_MODE` is set to "release" for production settings.
