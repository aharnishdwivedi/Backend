# Docker Setup for AI-Powered Incident Triage Assistant

This document provides comprehensive instructions for running the AI-Powered Incident Triage Assistant using Docker.

## 🐳 Quick Start

### Prerequisites
- Docker installed and running
- Docker Compose installed
- OpenAI API key

### 1. Clone and Setup
```bash
git clone <repository-url>
cd incident-triage-assistant
cp env.example .env
```

### 2. Configure Environment
Edit `.env` file:
```env
OPENAI_API_KEY=your_openai_api_key_here
```

### 3. Run with Docker Compose
```bash
# Using the provided script
./scripts/docker-run.sh

# Or manually
docker-compose up --build -d
```

### 4. Access the Application
- **API**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/api/v1/health
- **Frontend**: Serve `frontend/index.html` on http://localhost:3000

## 📦 Docker Images

### Backend Image
- **Base**: Alpine Linux (minimal footprint)
- **Runtime**: Go binary (statically compiled)
- **Size**: ~15MB (optimized multi-stage build)
- **Security**: Non-root user execution

### Database Image
- **Base**: MySQL 8.0 official image
- **Persistent Storage**: Docker volumes
- **Networking**: Isolated network

## 🔧 Docker Commands

### Development
```bash
# Build and start all services
make docker-dev

# View logs
make docker-logs

# Stop services
make docker-down

# Clean up everything
make docker-clean
```

### Manual Commands
```bash
# Build backend image
docker build -t incident-triage-assistant .

# Run backend only
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your_key \
  -e DB_HOST=your_db_host \
  incident-triage-assistant

# Run with docker-compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f backend
```

## 🏗️ Dockerfile Details

### Multi-Stage Build
```dockerfile
# Stage 1: Builder
FROM golang:1.21-alpine AS builder
# Compile the application

# Stage 2: Runtime
FROM alpine:latest
# Minimal runtime with only the binary
```

### Security Features
- Non-root user execution
- Minimal attack surface
- No shell access in production
- Health checks included

### Optimization
- Multi-stage build reduces image size
- Alpine Linux base for minimal footprint
- Static compilation for portability
- Layer caching for faster builds

## 🚀 Production Deployment

### 1. Build Production Image
```bash
./scripts/docker-build.sh v1.0.0
```

### 2. Push to Registry
```bash
docker tag incident-triage-assistant:latest your-registry/incident-triage-assistant:v1.0.0
docker push your-registry/incident-triage-assistant:v1.0.0
```

### 3. Deploy with Environment Variables
```bash
docker run -d \
  --name incident-triage-backend \
  -p 8080:8080 \
  -e OPENAI_API_KEY=your_production_key \
  -e DB_HOST=your_production_db \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=incident_triage \
  your-registry/incident-triage-assistant:v1.0.0
```

## 🔍 Monitoring and Health Checks

### Health Check Endpoint
```bash
curl http://localhost:8080/api/v1/health
```

### Docker Health Checks
- **Backend**: HTTP health check every 30s
- **Database**: MySQL ping every 20s
- **Startup**: Graceful startup with retries

### Logging
```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f mysql

# Follow logs with timestamps
docker-compose logs -f --timestamps
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Database Connection Issues
```bash
# Check if MySQL is running
docker-compose ps mysql

# Check MySQL logs
docker-compose logs mysql

# Connect to MySQL directly
docker-compose exec mysql mysql -u root -p
```

#### 2. Backend Startup Issues
```bash
# Check backend logs
docker-compose logs backend

# Check environment variables
docker-compose exec backend env | grep -E "(DB_|OPENAI_)"

# Restart backend
docker-compose restart backend
```

#### 3. Port Conflicts
```bash
# Check what's using port 8080
lsof -i :8080

# Use different ports
docker-compose up -d -p 8081:8080
```

### Debug Mode
```bash
# Run with debug logging
docker run -e LOG_LEVEL=debug incident-triage-assistant

# Access container shell (development only)
docker-compose exec backend sh
```

## 📊 Performance Optimization

### Resource Limits
```yaml
# In docker-compose.yml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
```

### Volume Mounts
```yaml
# For development
volumes:
  - ./logs:/app/logs
  - ./config:/app/config
```

## 🔒 Security Considerations

### Production Security
- Use secrets management for sensitive data
- Implement network policies
- Regular security updates
- Image scanning
- Non-root user execution

### Environment Variables
```bash
# Use Docker secrets in production
echo "your_openai_api_key" | docker secret create openai_api_key -

# Reference in docker-compose.yml
secrets:
  openai_api_key:
    external: true
```

## 📝 Docker Compose Configuration

### Services Overview
- **backend**: Go application with Echo framework
- **mysql**: MySQL 8.0 database
- **networks**: Isolated network for services
- **volumes**: Persistent database storage

### Environment Variables
- `OPENAI_API_KEY`: Required for AI functionality
- `DB_*`: Database connection parameters
- `SERVER_PORT`: Application port (default: 8080)

### Networking
- Services communicate via internal network
- External access through port mappings
- Health checks ensure service availability

## 🎯 Best Practices

### Development
1. Use `docker-compose` for local development
2. Mount source code for hot reloading
3. Use volume mounts for persistent data
4. Implement health checks

### Production
1. Use specific image tags (not `latest`)
2. Implement resource limits
3. Use secrets for sensitive data
4. Regular security updates
5. Monitor container health

### CI/CD
1. Build images in CI pipeline
2. Run tests in containers
3. Scan images for vulnerabilities
4. Deploy with rolling updates

---

For more information, see the main [README.md](README.md) file.
