#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_message() {
    local message=$1
    local color=$2
    local separator=$(printf '%*s' "${#message}" | tr ' ' '-')
    echo -e "${color}${separator}${NC}"
    echo -e "${color}${message}${NC}"
    echo -e "${color}${separator}${NC}"
}

print_message "Starting environment..." "${GREEN}"
source integration/.env && \
echo "GRPC_ADDRESS: $GRPC_ADDRESS" && \
docker-compose -f docker-compose-integration.yaml up -d --build
if [ $? -ne 0 ]; then
    print_message "Failed to start environment." "${RED}"
    docker-compose -f docker-compose-integration.yaml down
    exit 1
fi

print_message "Running tests..." "${GREEN}"
go test ./... -tags=integration
TEST_STATUS=$?

print_message "Cleaning up environment..." "${GREEN}"
docker-compose -f docker-compose-integration.yaml down

if [ $TEST_STATUS -eq 0 ]; then
    print_message "Tests passed successfully." "${GREEN}"
else
    print_message "Tests failed." "${RED}"
fi

exit $TEST_STATUS