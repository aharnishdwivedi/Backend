#!/bin/bash

# AI-Powered Incident Triage Assistant - Run Script

echo "🚀 Starting AI-Powered Incident Triage Assistant..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from template..."
    cp env.example .env
    echo "📝 Please edit .env file with your configuration before running again."
    echo "   Required: OPENAI_API_KEY, Database credentials"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if MySQL is running (basic check)
if ! command -v mysql &> /dev/null; then
    echo "⚠️  MySQL client not found. Please ensure MySQL is installed and running."
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Run tests
echo "🧪 Running tests..."
go test ./... -v

if [ $? -ne 0 ]; then
    echo "❌ Tests failed. Please fix the issues before running the application."
    exit 1
fi

echo "✅ Tests passed!"

# Run the application
echo "🚀 Starting the server..."
echo "📱 Frontend will be available at: http://localhost:3000"
echo "🔗 API will be available at: http://localhost:8080/api/v1"
echo "🏥 Health check: http://localhost:8080/api/v1/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run cmd/main.go
