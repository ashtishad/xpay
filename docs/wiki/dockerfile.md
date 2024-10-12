# Dockerfile Guide

## 1. Base Image Selection

Use specific version tags for base images. Provides a consistent and reproducible Go development environment.
```dockerfile
FROM golang:1.23.2-alpine3.20 AS builder
```
**Why**: Ensures all builds use the same Go version and Alpine Linux base, preventing inconsistencies across different build environments.

**Best Practices**:
- Use specific version tags to ensure reproducibility.
- Prefer Alpine-based images for smaller footprints.
- Regularly update base images to include security patches.

[Reference: Docker Official Images](https://docs.docker.com/develop/develop-images/baseimages/)

## 2. Working Directory

Set a working directory for your application. Establishes a dedicated workspace for the build process.
```dockerfile
WORKDIR /build
```
**Why**: Isolates build artifacts and prevents conflicts with system files.

**Best Practices**:
- Use meaningful directory names.
- Avoid using the root directory as the working directory.
- Be consistent with working directory paths across build stages.

[Reference: WORKDIR Instruction](https://docs.docker.com/engine/reference/builder/#workdir)

## 3. Dependency Management

Copy dependency files and install dependencies. Efficiently manages and verifies Go dependencies.
```dockerfile
COPY go.mod go.sum ./
RUN go mod download && go mod verify
```
**Why**: Ensures all required dependencies are available and verified before building, leveraging Docker's layer caching for faster subsequent builds.

**Best Practices**:
- Separate dependency installation from application code copying.
- Use `go mod verify` to ensure dependency integrity.
- Consider using a dependency caching mechanism for faster builds.

[Reference: Go Modules](https://go.dev/ref/mod#go-mod-verify)

## 4. Code Copy and Build

Copy application code and build the binary. Compiles the Go application into a lightweight, statically-linked binary.
```dockerfile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go
```
**Why**: Produces a self-contained executable optimized for size and compatibility with Alpine Linux.

**Best Practices**:
- Use build flags to optimize binary size.
- Disable CGO for better portability.
- Consider using multi-stage builds to keep the final image small.

[Reference: Go Build Flags](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

## 5. Multi-stage Build

Use the same Alpine version as the build stage for the final image. Creates a minimal runtime environment for the application.
```dockerfile
FROM alpine:3.20
```
**Why**: Significantly reduces the final image size by excluding build tools and intermediate files.

**Best Practices**:
- Use the smallest possible base image that supports your application.
- Consider using `scratch` for Go binaries if no additional runtime dependencies are needed.
- Keep the number of layers minimal in the final stage.

[Reference: Multi-stage Builds](https://docs.docker.com/develop/develop-images/multistage-build/)

## 6. Non-root User

Create and switch to a non-root user. Enhances container security by running the application as a non-privileged user.
```dockerfile
RUN adduser -D ash
USER ash
WORKDIR /home/ash
```
**Why**: Limits the potential impact of security vulnerabilities by restricting the application's permissions.

**Best Practices**:
- Always run applications as non-root users.
- Use the `--no-create-home` flag when it's not needed.
- Set appropriate file permissions for application files.

[Reference: Docker Security](https://docs.docker.com/engine/security/security/#linux-kernel-capabilities)

## 7. File Copying and Permissions

Copy built artifacts and set proper ownership. Transfers the compiled application and necessary files to the runtime environment with appropriate ownership.
```dockerfile
COPY --chown=ash:ash --from=builder /build .
COPY --chown=ash:ash --from=builder /build/migrations ./migrations
COPY --chown=ash:ash --from=builder /build/config.yaml .
```
**Why**: Ensures the application has the correct files and permissions to run properly as the non-root user.

**Best Practices**:
- Use `--chown` to set proper ownership in one step.
- Copy only necessary files to the final image.
- Use `.dockerignore` to exclude unnecessary files from the build context.

[Reference: COPY Instruction](https://docs.docker.com/engine/reference/builder/#copy)

## 8. Expose Ports

Declare the ports your application uses. Specifies the network interfaces on which the application will listen.
```dockerfile
EXPOSE 8080
```
**Why**: Provides metadata about the application's networking requirements, facilitating proper container networking setup.

**Best Practices**:
- Only expose necessary ports.
- Use specific port numbers rather than ranges.
- Document exposed ports in application documentation.

[Reference: EXPOSE Instruction](https://docs.docker.com/engine/reference/builder/#expose)

## 9. Entrypoint and CMD

Set the entrypoint for your application. Specifies how to run the application when the container starts.
```dockerfile
ENTRYPOINT ["/home/ash/main"]
CMD []
```
**Why**: Defines the primary command to execute the application, allowing for additional runtime arguments if needed.

**Best Practices**:
- Use `ENTRYPOINT` for the main command and `CMD` for default arguments.
- Prefer exec form (`[]`) over shell form for `ENTRYPOINT` and `CMD`.
- Provide a way to override the entrypoint for debugging purposes.

[Reference: ENTRYPOINT](https://docs.docker.com/engine/reference/builder/#entrypoint)

## Production vs Staging Considerations

- **Production**:
  - Optimize for security, performance, and minimal attack surface.
  - Implementation: Use minimal base images, implement security measures, and optimize configurations.
  - Best Practices:
    - Implement health checks.
    - Use read-only file systems where possible.
    - Implement proper logging and monitoring.

- **Staging**:
  - Mirror production environment while allowing for debugging and monitoring.
  - Implementation: Include additional tools for diagnostics while maintaining similarity to production setup.
  - Best Practices:
    - Use the same base image as production.
    - Include debugging tools without adding them to the production image.
    - Implement feature flags for easier testing.

[Reference: Docker Production Best Practices](https://docs.docker.com/develop/dev-best-practices/)
