# Be Your Moments - Complete CI/CD Pipeline

This repository contains a complete CI/CD pipeline for the Be Your Moments microservices application, supporting both **production deployment** (K3s) and **development deployment** (Docker Compose).

## 🏗️ Architecture Overview

The application consists of 5 microservices:

- **photo-svc** (Port 8001) - Photo management service
- **user-svc** (Port 8003) - User management service  
- **transaction-svc** (Port 8005) - Transaction processing service
- **upload-svc** (Port 8002) - File upload service
- **notification-svc** (Port 8004) - Notification service

## 🚀 CI/CD Pipeline Features

### **Production Pipeline** (K3s)
- **Trigger**: Tags like `v1.0.0`, `v2.1.0`
- **Deployment**: K3s cluster
- **Features**: Full microservices architecture with Kubernetes

### **Development Pipeline** (Docker Compose)
- **Trigger**: Tags like `x.dev.1.0.0`, `x.dev.2.1.0`
- **Deployment**: Docker Compose on VPS
- **Features**: Simplified development environment

## 📁 Project Structure

```
├── .github/workflows/
│   ├── ci-cd.yml                 # Production pipeline (K3s)
│   └── ci-cd-dev.yml            # Development pipeline (Docker Compose)
├── scripts/
│   ├── deploy.sh                 # Production deployment
│   ├── setup.sh                  # Production setup
│   ├── rollback.sh               # Production rollback
│   ├── monitor.sh                # Production monitoring
│   ├── deploy-dev.sh             # Development deployment
│   ├── setup-dev.sh              # Development setup
│   ├── monitor-dev.sh            # Development monitoring
│   └── k8s/                      # Kubernetes manifests
├── photo-svc/                    # Photo service
├── user-svc/                     # User service
├── transaction-svc/              # Transaction service
├── upload-svc/                   # Upload service
├── notification-svc/             # Notification service
└── pb/                          # Protocol buffers
```

## 🛠️ Quick Start

### Prerequisites

1. **VPS with Docker installed**
2. **GitHub repository with secrets configured**
3. **kubectl configured** (for production)

### 1. Production Setup (K3s)

```bash
# Setup K3s on VPS
./scripts/k3s-setup.sh

# Setup infrastructure
./scripts/setup.sh

# Create secrets
./scripts/create-secrets.sh

# Deploy services
./scripts/deploy.sh latest
```

### 2. Development Setup (Docker Compose)

```bash
# Setup development environment
./scripts/setup-dev.sh

# Start services
cd /tmp/be-yourmoments-dev
./start.sh
```

## 🔄 Deployment Process

### Production Deployment

1. **Create and push a production tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Run tests
   - Build Docker images
   - Deploy to K3s cluster

### Development Deployment

1. **Create and push a development tag:**
   ```bash
   git tag x.dev.1.0.0
   git push origin x.dev.1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Run tests
   - Build Docker images
   - Deploy to VPS using Docker Compose

## 📊 Monitoring

### Production Monitoring (K3s)

```bash
# Check status
./scripts/monitor.sh status

# Check health
./scripts/monitor.sh health

# View logs
./scripts/monitor.sh logs photo-svc
```

### Development Monitoring (Docker Compose)

```bash
# Check status
./scripts/monitor-dev.sh status

# Check health
./scripts/monitor-dev.sh health

# View logs
./scripts/monitor-dev.sh logs photo-svc
```

## 🔄 Rollback

### Production Rollback

```bash
# Rollback all services
./scripts/rollback.sh

# Rollback specific service
./scripts/rollback.sh photo-svc
```

### Development Rollback

```bash
# Restart services
cd /tmp/be-yourmoments-dev
./restart.sh

# Or use monitoring script
./scripts/monitor-dev.sh restart photo-svc
```

## 🏥 Health Checks

All services include health check endpoints:

- **photo-svc**: `http://localhost:8001/health`
- **user-svc**: `http://localhost:8003/health`
- **transaction-svc**: `http://localhost:8005/health`
- **upload-svc**: `http://localhost:8002/health`
- **notification-svc**: `http://localhost:8004/health`

## 🔧 Troubleshooting

### Production Issues

```bash
# Check pod status
kubectl get pods -n be-yourmoments

# Check logs
kubectl logs deployment/photo-svc -n be-yourmoments

# Check events
kubectl get events -n be-yourmoments
```

### Development Issues

```bash
# Check container status
docker ps

# Check logs
./logs.sh

# Check specific service
./scripts/monitor-dev.sh logs photo-svc
```

## 📈 Resource Requirements

### Production (K3s)
- **Minimum**: 4 cores, 8GB RAM
- **Recommended**: 8 cores, 16GB RAM
- **Services**: 5 microservices + infrastructure

### Development (Docker Compose)
- **Minimum**: 2 cores, 4GB RAM
- **Recommended**: 4 cores, 8GB RAM
- **Services**: 5 microservices + infrastructure

## 🔐 Security

### Secrets Management
- **Production**: Kubernetes secrets
- **Development**: Environment files
- **Service Account Keys**: Secure storage required

### Network Security
- **Production**: Network policies, ingress controllers
- **Development**: Docker networks, port mapping

## 📝 Best Practices

### Production
1. **Use semantic versioning** for tags
2. **Monitor resource usage** regularly
3. **Implement proper secrets management**
4. **Use network policies** for security
5. **Regular backups** of data

### Development
1. **Use development tags** for testing
2. **Monitor resource usage** locally
3. **Keep service account keys secure**
4. **Regular cleanup** of resources
5. **Test locally** before deploying

## 🆘 Support

### Production Support
- Check K3s logs: `kubectl logs deployment/<service> -n be-yourmoments`
- Check events: `kubectl get events -n be-yourmoments`
- Monitor resources: `kubectl top pods -n be-yourmoments`

### Development Support
- Check Docker logs: `docker logs <container>`
- Check service status: `./status.sh`
- Monitor resources: `docker stats`

## 📄 Documentation

- **CI_CD_README.md** - Production pipeline documentation
- **DEVELOPMENT_README.md** - Development pipeline documentation
- **scripts/README.md** - Script documentation
- **TROUBLESHOOTING.md** - Troubleshooting guide

## 🎯 Next Steps

### Production
1. **Setup K3s**: `./scripts/k3s-setup.sh`
2. **Deploy infrastructure**: `./scripts/setup.sh`
3. **Create secrets**: `./scripts/create-secrets.sh`
4. **Deploy services**: `./scripts/deploy.sh latest`
5. **Monitor**: `./scripts/monitor.sh status`

### Development
1. **Setup development**: `./scripts/setup-dev.sh`
2. **Start services**: `cd /tmp/be-yourmoments-dev && ./start.sh`
3. **Monitor**: `./scripts/monitor-dev.sh status`
4. **Create tag**: `git tag x.dev.1.0.0`
5. **Deploy**: `git push origin x.dev.1.0.0`

## 📊 Comparison

| Feature | Production (K3s) | Development (Docker Compose) |
|---------|------------------|------------------------------|
| **Complexity** | High | Low |
| **Resource Usage** | High | Low |
| **Scalability** | Excellent | Limited |
| **Management** | Complex | Simple |
| **Use Case** | Production | Development/Testing |
| **Setup Time** | 30+ minutes | 5-10 minutes |
| **Resource Requirements** | 4C/8GB+ | 2C/4GB+ |

## 🚀 Getting Started

Choose your deployment strategy:

### **For Production Use:**
- Use K3s pipeline with tags like `v1.0.0`
- Follow production documentation
- Use monitoring and rollback features

### **For Development/Testing:**
- Use Docker Compose pipeline with tags like `x.dev.1.0.0`
- Follow development documentation
- Use simplified management commands

Both pipelines are fully automated and will handle the complete deployment process from code to running services.
