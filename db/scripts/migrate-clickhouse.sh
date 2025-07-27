#!/bin/bash

# ClickHouse Migration Script for Sutram Trading Platform

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ACTION="up"
STEPS=1
VERSION=""
BINARY_PATH="./bin/ch-migrate"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "ClickHouse Migration Management for Sutram Trading Platform"
    echo ""
    echo "Options:"
    echo "  -a, --action ACTION     Migration action (up, down, steps, force, version, validate)"
    echo "  -s, --steps STEPS       Number of steps for 'steps' action (default: 1)"
    echo "  -v, --version VERSION   Version number for 'force' action"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --action=up                    # Run all pending ClickHouse migrations"
    echo "  $0 --action=down                  # Rollback all ClickHouse migrations"
    echo "  $0 --action=steps --steps=2       # Run 2 ClickHouse migration steps"
    echo "  $0 --action=steps --steps=-1      # Rollback 1 ClickHouse migration"
    echo "  $0 --action=force --version=5     # Force ClickHouse migration to version 5"
    echo "  $0 --action=version               # Show current ClickHouse migration version"
    echo "  $0 --action=validate              # Validate all ClickHouse migrations"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -a|--action)
            ACTION="$2"
            shift 2
            ;;
        -s|--steps)
            STEPS="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Build the migration binary if it doesn't exist
if [ ! -f "$BINARY_PATH" ]; then
    print_info "Building ClickHouse migration binary..."
    go build -o "$BINARY_PATH" cmd/ch-migrate/main.go
    if [ $? -ne 0 ]; then
        print_error "Failed to build ClickHouse migration binary"
        exit 1
    fi
    print_success "ClickHouse migration binary built successfully"
fi

# Construct the command
CMD="$BINARY_PATH -action=$ACTION"

if [ "$ACTION" = "steps" ]; then
    CMD="$CMD -steps=$STEPS"
elif [ "$ACTION" = "force" ]; then
    if [ -z "$VERSION" ]; then
        print_error "Version must be specified for force action"
        exit 1
    fi
    CMD="$CMD -version=$VERSION"
fi

# Run the migration command
print_info "Running ClickHouse migration: $ACTION"
print_info "Command: $CMD"

$CMD

if [ $? -eq 0 ]; then
    print_success "ClickHouse migration completed successfully"
else
    print_error "ClickHouse migration failed"
    exit 1
fi
