#!/bin/bash

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GATEWAY_NAME="sutram-gateway-temp"

# Check if the container is already running
if docker ps -q -f name=$GATEWAY_NAME | grep -q .; then
    echo "Stopping existing Nginx gateway container..."
    docker stop $GATEWAY_NAME
    docker rm $GATEWAY_NAME
fi

# Run Nginx in a Docker container
echo "Starting Nginx gateway on port 8010..."
docker run --name $GATEWAY_NAME \
  -p 8010:8010 \
  -v "$SCRIPT_DIR/nginx.conf:/etc/nginx/nginx.conf:ro" \
  -v "$SCRIPT_DIR/conf.d:/etc/nginx/conf.d:ro" \
  --add-host=host.docker.internal:host-gateway \
  -d nginx:alpine

# Check if container started successfully
if [ $? -ne 0 ]; then
    echo "Failed to start Nginx gateway container."
    exit 1
fi

echo "Nginx gateway is running. Press Ctrl+C to stop."
echo "Access the gateway at http://localhost:8010"

# Wait for user interrupt
trap "docker stop $GATEWAY_NAME; docker rm $GATEWAY_NAME; echo 'Nginx stopped.'" INT TERM
read -r -d '' _ </dev/tty

# Clean up shouldn't reach here, but just in case
docker stop $GATEWAY_NAME
docker rm $GATEWAY_NAME
echo "Nginx gateway has been stopped."
