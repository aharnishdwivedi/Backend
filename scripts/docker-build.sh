#!/bin/bash

# AI-Powered Incident Triage Assistant - Docker Build Script

echo "🔨 Building AI-Powered Incident Triage Assistant Docker image..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

# Set image name and tag
IMAGE_NAME="incident-triage-assistant"
TAG=${1:-latest}

echo "📦 Building image: $IMAGE_NAME:$TAG"

# Build the Docker image
docker build -t $IMAGE_NAME:$TAG .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully!"
    echo "📋 Image details:"
    docker images $IMAGE_NAME:$TAG
    
    echo ""
    echo "🚀 To run the container:"
    echo "   docker run -p 8080:8080 -e OPENAI_API_KEY=your_key -e DB_HOST=your_db_host $IMAGE_NAME:$TAG"
    echo ""
    echo "📋 To push to registry:"
    echo "   docker tag $IMAGE_NAME:$TAG your-registry/$IMAGE_NAME:$TAG"
    echo "   docker push your-registry/$IMAGE_NAME:$TAG"
else
    echo "❌ Docker build failed!"
    exit 1
fi
