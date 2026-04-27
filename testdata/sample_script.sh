#!/bin/bash

# Sample shell script for testing clipper capabilities
# This script demonstrates various shell scripting patterns

set -e  # Exit on any error

# Configuration variables
SCRIPT_NAME="$(basename "$0")"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/tmp/${SCRIPT_NAME}.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    local level="$1"
    local message="$2"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [$level] $message" >> "$LOG_FILE"
    echo -e "${GREEN}[$level]${NC} $message"
}

# Error logging
error() {
    local message="$1"
    echo -e "${RED}[ERROR]${NC} $message" >&2
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] $message" >> "$LOG_FILE"
}

# Warning logging
warning() {
    local message="$1"
    echo -e "${YELLOW}[WARNING]${NC} $message"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [WARNING] $message" >> "$LOG_FILE"
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Validate input
validate_input() {
    local input="$1"
    if [[ -z "$input" ]]; then
        error "Input cannot be empty"
        return 1
    fi
    return 0
}

# Process files in a directory
process_files() {
    local directory="$1"
    local pattern="$2"

    if [[ ! -d "$directory" ]]; then
        error "Directory does not exist: $directory"
        return 1
    fi

    log "INFO" "Processing files in $directory with pattern $pattern"

    local count=0
    while IFS= read -r -d '' file; do
        log "INFO" "Processing file: $file"
        # Simulate file processing
        sleep 0.1
        ((count++))
    done < <(find "$directory" -name "$pattern" -type f -print0)

    log "INFO" "Processed $count files"
}

# Main function
main() {
    log "INFO" "Starting $SCRIPT_NAME"

    # Check dependencies
    if ! command_exists "find"; then
        error "Required command 'find' not found"
        exit 1
    fi

    # Parse command line arguments
    local target_dir=""
    local file_pattern="*.txt"

    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--directory)
                target_dir="$2"
                shift 2
                ;;
            -p|--pattern)
                file_pattern="$2"
                shift 2
                ;;
            -h|--help)
                echo "Usage: $SCRIPT_NAME [OPTIONS]"
                echo "Options:"
                echo "  -d, --directory DIR    Target directory (default: current)"
                echo "  -p, --pattern PATTERN  File pattern (default: *.txt)"
                echo "  -h, --help            Show this help"
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                exit 1
                ;;
        esac
    done

    # Set default directory
    if [[ -z "$target_dir" ]]; then
        target_dir="$(pwd)"
    fi

    # Validate and process
    if validate_input "$target_dir"; then
        process_files "$target_dir" "$file_pattern"
    fi

    log "INFO" "$SCRIPT_NAME completed successfully"
}

# Cleanup function
cleanup() {
    log "INFO" "Cleaning up"
    # Add cleanup logic here
}

# Set trap for cleanup
trap cleanup EXIT

# Run main function with all arguments
main "$@"