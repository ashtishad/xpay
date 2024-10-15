#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

cp "$PROJECT_ROOT/local-dev/Makefile.dev" "$PROJECT_ROOT/Makefile"
cp "$PROJECT_ROOT/local-dev/compose.yaml.dev" "$PROJECT_ROOT/compose.yaml"

if [ ! -f "$PROJECT_ROOT/config.yaml" ]; then
    cp "$PROJECT_ROOT/local-dev/config.yaml.example" "$PROJECT_ROOT/config.yaml"
fi

# Function to install a Go package if it's not already installed
install_if_not_exists() {
    if ! command -v $1 &> /dev/null; then
        echo "Installing $1..."
        go install $2
    else
        echo "$1 is already installed."
    fi
}

# Install development dependencies if not already installed
install_if_not_exists air github.com/cosmtrek/air@latest
install_if_not_exists golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint@latest
install_if_not_exists swag github.com/swaggo/swag/cmd/swag@latest
install_if_not_exists migrate "github.com/golang-migrate/migrate/v4/cmd/migrate@latest" "-tags 'postgres'"

echo "Dev environment set up complete!"
