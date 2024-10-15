#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

cp "$PROJECT_ROOT/local-dev/Makefile.prod" "$PROJECT_ROOT/Makefile"
cp "$PROJECT_ROOT/local-dev/compose.yaml.prod" "$PROJECT_ROOT/compose.yaml"

if [ ! -f "$PROJECT_ROOT/config.yaml" ]; then
    cp "$PROJECT_ROOT/local-dev/config.yaml.example" "$PROJECT_ROOT/config.yaml"
fi

echo "Production environment set up complete!"
