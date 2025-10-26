# Be Your Moments - Development Environment

This document describes the development environment setup for Be Your Moments microservices using Docker Compose.

## 🏗️ Architecture Overview

The development environment consists of:

- **5 Microservices** (photo-svc, user-svc, transaction-svc, upload-svc, notification-svc)
- **Infrastructure Services** (PostgreSQL, Redis, NATS, Consul, MinIO)
- **Docker Compose** for local development

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Git repository cloned
- Service account keys (for user-svc and notification-svc)

### 1. Setup Development Environment

```bash
# Run setup script
./scripts/setup-dev.sh
```

### 2. Start Services

```bash
# Navigate to development directory
cd /tmp/be-yourmoments-dev

# Start all services
./start.sh

# Or manually
docker-compose -f docker-compose-dev.yaml up -d
```

### 3. Check Status

```bash
# Check service status
./status.sh

# Or use monitoring script
./scripts/monitor-dev.sh status
```

## 📁 Development Directory Structure

```
/tmp/be-yourmoments-dev/
├── docker-compose-dev.yaml    # Development compose file
├── .env                      # Infrastructure environment
├── .env.photo                # Photo service environment
├── .env.user                 # User service environment
├── .env.transaction          # Transaction service environment
├── .env.upload               # Upload service environment
├── .env.notification         # Notification service environment
├── start.sh                  # Start all services
├── stop.sh                   # Stop all services
├── restart.sh                # Restart all services
├── status.sh                 # Show service status
├── logs.sh                   # Show service logs
├── init/                     # Database initialization scripts
└── scripts/                  # Management scripts
```

## 🔧 Management Commands

### Basic Operations

```bash
# Start services
./start.sh

# Stop services
./stop.sh

# Restart services
./restart.sh

# Show status
./status.sh

# Show logs
./logs.sh
```

### Advanced Operations

```bash
# Monitor services
./scripts/monitor-dev.sh status

# Check health
./scripts/monitor-dev.sh health

# Follow logs for specific service
./scripts/monitor-dev.sh follow photo-svc

# Show resource usage
./scripts/monitor-dev.sh resources

# Execute command in container
./scripts/monitor-dev.sh exec photo-svc bash
```

## 🌐 Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Photo Service | http://localhost:8001 | Photo management API |
| User Service | http://localhost:8003 | User management API |
| Transaction Service | http://localhost:8005 | Transaction processing API |
| Upload Service | http://localhost:8002 | File upload API |
| Notification Service | http://localhost:8004 | Notification API |
| MinIO Console | http://localhost:9001 | Object storage console |
| Consul UI | http://localhost:8500 | Service discovery UI |

## 🔍 Monitoring and Debugging

### Health Checks

```bash
# Check all service health
./scripts/monitor-dev.sh health

# Check specific service
curl http://localhost:8001/health
```

### Logs

```bash
# Show logs for all services
./logs.sh

# Show logs for specific service
./scripts/monitor-dev.sh logs photo-svc

# Follow logs in real-time
./scripts/monitor-dev.sh follow user-svc
```

### Resource Usage

```bash
# Show resource usage
./scripts/monitor-dev.sh resources

# Show detailed status
./scripts/monitor-dev.sh all
```

## 🗄️ Database Management

### PostgreSQL

```bash
# Connect to database
docker exec -it be-yourmoments-dev_postgres_1 psql -U postgres

# Run migrations
docker exec -it be-yourmoments-dev_photo-svc_1 ./migrate up
```

### Redis

```bash
# Connect to Redis
docker exec -it be-yourmoments-dev_redis_1 redis-cli

# Monitor Redis
docker exec -it be-yourmoments-dev_redis_1 redis-cli monitor
```

## 🔐 Configuration

### Environment Variables

Each service has its own environment file:

- `.env.photo` - Photo service configuration
- `.env.user` - User service configuration
- `.env.transaction` - Transaction service configuration
- `.env.upload` - Upload service configuration
- `.env.notification` - Notification service configuration

### Service Account Keys

For services that require Google Cloud integration:

1. **User Service**: Place service account key in `user-svc/serviceAccountKey.json`
2. **Notification Service**: Place service account key in `notification-svc/serviceAccountKey.json`

## 🚀 CI/CD Development Pipeline

### GitHub Actions Workflow

The development pipeline is triggered by tags starting with `x.dev`:

```bash
# Create development tag
git tag x.dev.1.0.0
git push origin x.dev.1.0.0
```

### Workflow Features

- **Automatic testing** on pull requests
- **Docker image building** on tag push
- **Automatic deployment** to VPS using Docker Compose
- **Environment setup** with proper configuration

## 🔧 Troubleshooting

### Common Issues

1. **Services not starting:**
   ```bash
   # Check Docker status
   docker ps
   
   # Check logs
   ./logs.sh
   
   # Restart services
   ./restart.sh
   ```

2. **Database connection issues:**
   ```bash
   # Check PostgreSQL status
   docker exec -it be-yourmoments-dev_postgres_1 pg_isready
   
   # Check database logs
   docker logs be-yourmoments-dev_postgres_1
   ```

3. **Service discovery issues:**
   ```bash
   # Check Consul status
   curl http://localhost:8500/v1/status/leader
   
   # Check service registration
   curl http://localhost:8500/v1/catalog/services
   ```

### Debugging Commands

```bash
# Show all containers
docker ps -a

# Show container logs
docker logs <container_name>

# Execute command in container
docker exec -it <container_name> bash

# Show network information
docker network ls

# Show volume information
docker volume ls
```

## 📊 Performance Monitoring

### Resource Usage

```bash
# Show resource usage
./scripts/monitor-dev.sh resources

# Show detailed statistics
docker stats
```

### Service Health

```bash
# Check service health
./scripts/monitor-dev.sh health

# Check specific service
curl http://localhost:8001/health
```

## 🔄 Development Workflow

### 1. Code Changes

```bash
# Make code changes
# Test locally
go test ./...

# Build and test
docker build -t test-photo-svc ./photo-svc
```

### 2. Deploy Changes

```bash
# Create development tag
git tag x.dev.1.0.1
git push origin x.dev.1.0.1

# GitHub Actions will automatically deploy
```

### 3. Monitor Deployment

```bash
# Check deployment status
./scripts/monitor-dev.sh status

# Check logs
./scripts/monitor-dev.sh logs photo-svc
```

## 🧹 Cleanup

### Stop Services

```bash
# Stop all services
./stop.sh

# Remove containers and volumes
docker-compose -f docker-compose-dev.yaml down -v
```

### Cleanup Resources

```bash
# Remove unused images
docker image prune -f

# Remove unused volumes
docker volume prune -f

# Remove unused networks
docker network prune -f
```

## 📝 Best Practices

1. **Always use development tags** for testing
2. **Monitor resource usage** regularly
3. **Check service health** before deploying
4. **Use proper environment variables**
5. **Keep service account keys secure**
6. **Regular cleanup** of unused resources

## 🆘 Support

For issues and questions:

1. Check the logs: `./logs.sh`
2. Check service status: `./status.sh`
3. Check health: `./scripts/monitor-dev.sh health`
4. Review resource usage: `./scripts/monitor-dev.sh resources`

## 📄 Next Steps

1. **Setup development environment**: `./scripts/setup-dev.sh`
2. **Start services**: `cd /tmp/be-yourmoments-dev && ./start.sh`
3. **Monitor services**: `./scripts/monitor-dev.sh status`
4. **Create development tag**: `git tag x.dev.1.0.0`
5. **Deploy to VPS**: `git push origin x.dev.1.0.0`
