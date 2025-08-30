#!/bin/bash

# AI-Powered Incident Triage Assistant - Docker Run Script

echo "🐳 Starting AI-Powered Incident Triage Assistant with Docker..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker and Docker Compose."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from template..."
    cp env.example .env
    echo "📝 Please edit .env file with your OPENAI_API_KEY before running again."
    echo "   Required: OPENAI_API_KEY=your_openai_api_key_here"
    exit 1
fi

# Check if OPENAI_API_KEY is set
if ! grep -q "OPENAI_API_KEY=" .env || grep -q "OPENAI_API_KEY=your_openai_api_key_here" .env; then
    echo "❌ OPENAI_API_KEY not configured in .env file."
    echo "   Please set your OpenAI API key in the .env file:"
    echo "   OPENAI_API_KEY=your_actual_api_key_here"
    exit 1
fi

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker-compose down

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

# Check backend health
echo "🏥 Checking backend health..."
for i in {1..30}; do
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo "✅ Backend is healthy!"
        break
    fi
    echo "⏳ Waiting for backend to be ready... (attempt $i/30)"
    sleep 2
done

if [ $i -eq 30 ]; then
    echo "❌ Backend failed to start properly. Check logs with: docker-compose logs backend"
    exit 1
fi

echo ""
echo "🎉 AI-Powered Incident Triage Assistant is running!"
echo ""
echo "📱 Frontend: http://localhost:3000 (serve frontend/index.html)"
echo "🔗 API: http://localhost:8080/api/v1"
echo "🏥 Health Check: http://localhost:8080/api/v1/health"
echo "🗄️  Database: localhost:3306 (incident_triage)"
echo ""
echo "📋 Useful Commands:"
echo "   View logs: docker-compose logs -f"
echo "   Stop services: docker-compose down"
echo "   Restart: docker-compose restart"
echo "   Clean up: docker-compose down -v"
echo ""
echo "Press Ctrl+C to stop the services"
echo ""

# Show logs
docker-compose logs -f
