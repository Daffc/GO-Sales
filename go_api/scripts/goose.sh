#!/bin/bash

# Sourcing local .env file.
source .env

# Calling goose with connection arguments
GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose "$@"