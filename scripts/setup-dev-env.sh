#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

cp "$PROJECT_ROOT/env-configs/dev.Makefile" "$PROJECT_ROOT/Makefile"
cp "$PROJECT_ROOT/env-configs/compose.dev.yaml" "$PROJECT_ROOT/compose.yaml"

if [ ! -f "$PROJECT_ROOT/config.example.yaml" ]; then
    cp "$PROJECT_ROOT/env-configs/config.example.yaml" "$PROJECT_ROOT/config.example.yaml"
fi

# Add GOPATH/bin to PATH if not already present
GOPATH=$(go env GOPATH)
if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
    export PATH="$PATH:$GOPATH/bin"

    # Add to shell profile if not already there
    SHELL_PROFILE=""
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_PROFILE="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        SHELL_PROFILE="$HOME/.bashrc"
    fi

    if [ -n "$SHELL_PROFILE" ]; then
        if ! grep -q "export PATH=\$PATH:\$GOPATH/bin" "$SHELL_PROFILE"; then
            echo 'export PATH=$PATH:$GOPATH/bin' >> "$SHELL_PROFILE"
            echo "Added GOPATH/bin to $SHELL_PROFILE"
        fi
    fi
fi

# Function to install a Go package if it's not already installed
install_if_not_exists() {
    if ! command -v $1 &> /dev/null; then
        echo "Installing $1..."
        go install $2

        # Verify installation
        if ! command -v $1 &> /dev/null; then
            echo "Warning: $1 installation may have failed. Please check your GOPATH and PATH settings."
            echo "Current GOPATH: $GOPATH"
            echo "Current PATH: $PATH"
        else
            echo "$1 installed successfully."
        fi
    else
        echo "$1 is already installed."
    fi
}

# Install development dependencies if not already installed
install_if_not_exists air github.com/air-verse/air@latest
install_if_not_exists golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64
install_if_not_exists swag github.com/swaggo/swag/cmd/swag@latest
install_if_not_exists migrate "github.com/golang-migrate/migrate/v4/cmd/migrate@latest" "-tags 'postgres'"

# Setup Git Hooks
setup_git_hooks() {
    mkdir -p "$PROJECT_ROOT/.git/hooks"
    cp "$PROJECT_ROOT/scripts/pre-push" "$PROJECT_ROOT/.git/hooks/pre-push"
    chmod +x "$PROJECT_ROOT/.git/hooks/pre-push"
    echo "Git hooks setup complete."
}

setup_git_hooks

# Final verification
echo "Verifying installations:"
for tool in air golangci-lint swag migrate; do
    if command -v $tool &> /dev/null; then
        echo "✓ $tool is available at $(which $tool)"
    else
        echo "✗ $tool is not available in PATH"
    fi
done

echo "Dev environment set up complete!"
echo "Note: You may need to restart your terminal or run 'source $SHELL_PROFILE' for PATH changes to take effect."
