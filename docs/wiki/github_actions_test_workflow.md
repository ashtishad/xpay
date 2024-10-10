# GitHub Actions CI Guide

This guide explains our Continuous Integration (CI) workflow using GitHub Actions.

## Workflow Overview

Our CI pipeline automatically tests and lints the codebase on every push and pull request to the main branch, specifically when changes are made to Go files, SQL files, Go module files, or CI workflow files.

## Workflow Configuration

The workflow is defined in `.github/workflows/test.yaml`. Here's a step-by-step breakdown:

### 1. Workflow Trigger

```yaml
on:
  push:
    branches: [ "main" ]
    paths:
      - '**.go'
      - '**.sql'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/**.yaml'
  pull_request:
    branches: [ "main" ]
    paths:
      - '**.go'
      - '**.sql'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/**.yaml'
```

This section defines when the workflow runs. It triggers on pushes and pull requests to the main branch, but only when specific files are changed. This optimizes CI resources by running only when relevant files are modified.

### 2. Environment Variables

```yaml
env:
  GO_VERSION: '1.23.2'
  GOLANGCI_LINT_VERSION: v1.61
```

Environment variables are set here for easy management and consistency across the workflow.

### 3. Job Setup

```yaml
jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
```

This defines a job named "Test" that runs on the latest Ubuntu runner, ensuring access to recent packages and tools.

### 4. PostgreSQL Service

```yaml
services:
  postgres:
    image: postgres:17.0-alpine3.20
    env:
      POSTGRES_DB: xpay
      POSTGRES_USER: ash
      POSTGRES_PASSWORD: samplepass
    ports:
      - "5432:5432"
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```

This sets up a PostgreSQL service container for testing. The health check ensures the database is ready before proceeding.

### 5. Workflow Steps

```yaml
steps:
- uses: actions/checkout@v4

- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}

- name: Update Go modules
  run: go mod tidy

- name: Run migrations
  run: make migrate-up

- name: Run Tests
  run: make test

- name: Golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: ${{ env.GOLANGCI_LINT_VERSION }}
```

These steps:
1. Check out the code
2. Set up the Go environment
3. Update Go modules
4. Run database migrations
5. Run tests
6. Perform linting

## Best Practices

1. **Selective Triggering**: Run CI only on relevant file changes.
2. **Version Consistency**: Use environment variables for version numbers.
3. **Database Management**: Use a service container with health checks.
4. **Comprehensive Testing**: Include running migrations, tests, and linting.
5. **Up-to-date Dependencies**: Regularly update Go modules.
6. **Use of Makefile**: Utilize make commands for common operations like migrations and testing.

### Debugging CI Issues

If you encounter issues in the CI environment:

1. **Check Logs**: Examine the full GitHub Actions logs for error messages and stack traces.

2. **Environment Differences**:
   - CI uses PostgreSQL on `127.0.0.1:5432`. Ensure your app can connect to this address.
   - Verify that all required environment variables are set in the workflow file.

3. **File Paths**:
   - Ensure all relative paths in your code and Makefile are correct.

4. **Database Migrations**:
   - Check if migrations are running successfully before tests.

5. **Reproduce Locally**:
   Try to reproduce the CI environment locally using Docker or by running the make commands directly.

6. **Add Debug Steps**:
   If needed, temporarily add debug steps to your workflow:

   ```yaml
   - name: Debug Info
     run: |
       pwd
       ls -la
       go version
       psql --version
   ```

7. **Check Action Versions**:
   Ensure you're using the latest stable versions of GitHub Actions.

## Modifying the Workflow

1. Edit `.github/workflows/test.yaml`.
2. Commit and push changes to trigger a new CI run.
3. Monitor the Actions tab in your GitHub repository for results.

Remember to keep sensitive information, such as database credentials, secure by using GitHub Secrets for production environments.
