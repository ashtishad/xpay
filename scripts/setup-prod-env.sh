#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

cp "$PROJECT_ROOT/env-configs/prod.Makefile" "$PROJECT_ROOT/Makefile"
cp "$PROJECT_ROOT/env-configs/compose.prod.yaml" "$PROJECT_ROOT/compose.yaml"

if [ ! -f "$PROJECT_ROOT/config.example.yaml" ]; then
    cp "$PROJECT_ROOT/env-configs/config.example.yaml" "$PROJECT_ROOT/config.example.yaml"
fi

# Setup Git Hooks
setup_git_hooks() {
    mkdir -p "$PROJECT_ROOT/.git/hooks"
    cp "$PROJECT_ROOT/scripts/pre-push" "$PROJECT_ROOT/.git/hooks/pre-push"
    chmod +x "$PROJECT_ROOT/.git/hooks/pre-push"
    echo "Git hooks setup complete."
}

setup_git_hooks

echo "Production environment set up complete!"
