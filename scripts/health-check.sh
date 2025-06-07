#!/bin/bash

# Health Check Script for URL Shortener API
# Usage: ./scripts/health-check.sh [URL]

set -e

# Configuration
URL=${1:-"http://localhost:8080"}
TIMEOUT=10
RETRIES=3

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_endpoint() {
    local endpoint=$1
    local expected_status=${2:-200}
    local description=$3
    
    echo -n "Checking $description... "
    
    for i in $(seq 1 $RETRIES); do
        if response=$(curl -s -w "%{http_code}" --max-time $TIMEOUT "$URL$endpoint" 2>/dev/null); then
            status_code="${response: -3}"
            if [ "$status_code" = "$expected_status" ]; then
                echo -e "${GREEN}✓${NC}"
                return 0
            fi
        fi
        
        if [ $i -lt $RETRIES ]; then
            echo -n "."
            sleep 2
        fi
    done
    
    echo -e "${RED}✗${NC} (Status: ${status_code:-'No response'})"
    return 1
}

# Main health check
main() {
    log_info "Starting health check for URL Shortener API"
    log_info "Target URL: $URL"
    echo

    # Check if URL is reachable
    if ! curl -s --max-time $TIMEOUT "$URL" >/dev/null 2>&1; then
        log_error "Cannot reach $URL"
        exit 1
    fi

    # Health checks
    local failed=0

    # Basic health endpoint
    if ! check_endpoint "/health" "200" "Health endpoint"; then
        ((failed++))
    fi

    # Stats endpoint
    if ! check_endpoint "/stats" "200" "Stats endpoint"; then
        ((failed++))
    fi

    # Swagger documentation
    if ! check_endpoint "/swagger/index.html" "200" "Swagger UI"; then
        ((failed++))
    fi

    # API documentation JSON
    if ! check_endpoint "/swagger/doc.json" "200" "API documentation"; then
        ((failed++))
    fi

    echo

    # Summary
    if [ $failed -eq 0 ]; then
        log_info "All health checks passed! ✅"
        
        # Additional info
        echo
        log_info "Service Information:"
        if health_data=$(curl -s --max-time $TIMEOUT "$URL/health" 2>/dev/null); then
            echo "$health_data" | python3 -m json.tool 2>/dev/null || echo "$health_data"
        fi
        
        echo
        log_info "Statistics:"
        if stats_data=$(curl -s --max-time $TIMEOUT "$URL/stats" 2>/dev/null); then
            echo "$stats_data" | python3 -m json.tool 2>/dev/null || echo "$stats_data"
        fi
        
        exit 0
    else
        log_error "$failed health check(s) failed! ❌"
        exit 1
    fi
}

# Help function
show_help() {
    echo "URL Shortener Health Check Script"
    echo
    echo "Usage: $0 [URL]"
    echo
    echo "Arguments:"
    echo "  URL    Base URL of the API (default: http://localhost:8080)"
    echo
    echo "Examples:"
    echo "  $0                                    # Check localhost"
    echo "  $0 http://localhost:8080              # Check specific port"
    echo "  $0 https://your-domain.com            # Check production"
    echo
    echo "Exit codes:"
    echo "  0    All checks passed"
    echo "  1    One or more checks failed"
}

# Parse arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac
