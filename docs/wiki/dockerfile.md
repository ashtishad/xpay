# Dockerfile Guide

## 1. Base Image Selection

**How**: Use specific version tags for base images.
```dockerfile
FROM golang:1.23.2-alpine3.20 AS builder
```

**Why**: Ensures reproducibility and prevents unexpected changes from image updates.

**Best Practice**:
- Use slim or alpine variants for smaller image sizes.
- For production, consider distroless images for enhanced security.

[Reference: Docker Official Images](https://docs.docker.com/develop/develop-images/baseimages/)

## 2. Working Directory

**How**: Set a working directory for your application.
```dockerfile
WORKDIR /app
```

**Why**: Organizes files within the container and sets the context for subsequent commands.

**Best Practice**: Use a descriptive directory name relevant to your application.

[Reference: WORKDIR Instruction](https://docs.docker.com/engine/reference/builder/#workdir)

## 3. Dependency Management

**How**: Copy dependency files and install dependencies.
```dockerfile
COPY go.mod go.sum ./
RUN go mod download && go mod verify
```

**Why**: Leverages Docker's layer caching for faster builds when dependencies don't change.

**Best Practice**: Verify downloaded modules for integrity.

[Reference: Go Modules](https://go.dev/ref/mod#go-mod-verify)

## 4. Code Copy and Build

**How**: Copy application code and build the binary.
```dockerfile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o xm main.go
```

**Why**: Builds a statically-linked binary optimized for size.

**Best Practice**: Use build flags to reduce binary size and disable CGO for better portability.

[Reference: Go Build Flags](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

## 5. Multi-stage Build

**How**: Use the same Alpine version as the build stage for the final image.
```dockerfile
FROM alpine:3.20
```

**Why**: Reduces the final image size by excluding build tools and intermediate files.

**Best Practice**: Use the smallest possible base image that supports your application.

[Reference: Multi-stage Builds](https://docs.docker.com/develop/develop-images/multistage-build/)

## 6. Non-root User

**How**: Create and switch to a non-root user.
```dockerfile
RUN adduser -D ash
USER ash
WORKDIR /home/ash
```

**Why**: Enhances security by limiting potential damage from container breakouts.

**Best Practice**: Always run containers as non-root users in production.

[Reference: Docker Security](https://docs.docker.com/engine/security/security/#linux-kernel-capabilities)

## 7. File Copying and Permissions

**How**: Copy built artifacts and set proper ownership.
```dockerfile
COPY --chown=ash:ash --from=builder /app/ .
COPY --chown=ash:ash --from=builder /app/migrations ./migrations
COPY --chown=ash:ash --from=builder /app/config.yaml .
```

**Why**: Ensures the application files are owned by the non-root user.

**Best Practice**: Always set appropriate file permissions and ownership.

[Reference: COPY Instruction](https://docs.docker.com/engine/reference/builder/#copy)

## 8. Expose Ports

**How**: Declare the ports your application uses.
```dockerfile
EXPOSE 8080
```

**Why**: Documents the ports the application uses, improving clarity for operators.

**Best Practice**: Only expose necessary ports.

[Reference: EXPOSE Instruction](https://docs.docker.com/engine/reference/builder/#expose)

## 9. Entrypoint and CMD

**How**: Set the entrypoint for your application.
```dockerfile
ENTRYPOINT ["/home/ash/xm"]
CMD []
```

**Why**: Defines how your container will run as an executable.

**Best Practice**: Use ENTRYPOINT for the main command and CMD for default arguments.

[Reference: ENTRYPOINT](https://docs.docker.com/engine/reference/builder/#entrypoint)

## Production vs Staging Considerations

- **Production**:
  - Use minimal, security-hardened base images.
  - Implement strict security measures (non-root users, read-only filesystems).
  - Optimize for performance and minimal attack surface.

- **Staging**:
  - Can use larger images with debugging tools if needed.
  - May include additional monitoring or logging tools.
  - Should mirror production as closely as possible for accurate testing.

[Reference: Docker Production Best Practices](https://docs.docker.com/develop/dev-best-practices/)
