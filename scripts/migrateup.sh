#!/bin/bash

# From Learn CI/CD starter
ENV_PATH=".env"

# Load .env and export variables
if [ -f "$ENV_PATH" ]; then
    set -a
    source "$ENV_PATH"
    set +a
fi

# Sanity check
if [ -z "$DB_URL" ]; then
    echo "❌ DB_URL is not set"
    exit 1
fi

cd sql/schema
goose postgres "$DB_URL" up