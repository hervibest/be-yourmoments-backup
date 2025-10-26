# Be Your Moments - CI/CD Pipeline

This repository contains a complete CI/CD pipeline for the Be Your Moments microservices application, designed to work with GitHub Actions and K3s.

## 🏗️ Architecture Overview

The application consists of 5 microservices:

- **photo-svc** (Port 8001) - Photo management service
- **user-svc** (Port 8003) - User management service  
- **transaction-svc** (Port 8005) - Transaction processing service
- **upload-svc** (Port 8002) - File upload service
- **notification-svc** (Port 8004) - Notification service

## 🚀 CI/CD Pipeline Features

### GitHub Actions Workflow
- **Automatic testing** on pull requests
- **Docker image building** on tag push
- **Automatic deployment** to K3s cluster
- **Multi-service support** with parallel builds

### Deployment Process
1. **Test Stage** - Runs unit tests and linting
2. **Build Stage** - Builds Docker images for all services
3. **Deploy Stage** - Deploys to K3s cluster via SSH

## 📁 Project Structure

```
├── .github/workflows/
│   └── ci-cd.yml                 # GitHub Actions workflow
├── scripts/
│   ├── deploy.sh                 # Main deployment script
│   ├── setup.sh                  # Infrastructure setup
│   ├── rollback.sh               # Rollback deployments
│   ├── monitor.sh                # Monitoring and troubleshooting
│   ├── create-secrets.sh         # Create Kubernetes secrets
│   ├── setup-env.sh              # Environment setup
│   ├── k3s-setup.sh              # K3s installation script
│   ├── complete-setup.sh         # Complete setup script
│   └── k8s/                      # Kubernetes manifests
│       ├── photo-svc.yaml
│       ├── user-svc.yaml
│       ├── transaction-svc.yaml
│       ├── upload-svc.yaml
│       ├── notification-svc.yaml
│       ├── infrastructure.yaml
│       ├── secrets.yaml
│       └── ingress.yaml
├── photo-svc/                    # Photo service
├── user-svc/                     # User service
├── transaction-svc/              # Transaction service
├── upload-svc/                   # Upload service
├── notification-svc/             # Notification service
└── pb/                          # Protocol buffers
```

## 🛠️ Quick Start

### Prerequisites

1. **VPS with K3s installed**
2. **GitHub repository with secrets configured**
3. **kubectl configured to access K3s cluster**

### 1. Initial Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd be-yourmoments-backup

# Make scripts executable
chmod +x scripts/*.sh

# Run complete setup
./scripts/complete-setup.sh
```

### 2. Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions, and add:

```
SSH_PRIVATE_KEY    # Private key for SSH access to VPS
VPS_USER          # Username for VPS SSH access  
VPS_HOST          # IP address or hostname of your VPS
```

### 3. Setup K3s on VPS

```bash
# On your VPS, run:
./scripts/k3s-setup.sh
```

### 4. Deploy Infrastructure

```bash
# Setup infrastructure
./scripts/setup.sh

# Create secrets
./scripts/create-secrets.sh interactive

# Deploy services
./scripts/deploy.sh latest
```

### 5. Monitor Deployment

```bash
# Check status
./scripts/monitor.sh status

# Check health
./scripts/monitor.sh health

# View logs
./scripts/monitor.sh logs photo-svc
```

## 🔄 Deployment Process

### Automatic Deployment

1. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Run tests
   - Build Docker images
   - Push to GitHub Container Registry
   - Deploy to K3s cluster

### Manual Deployment

```bash
# Deploy specific version
./scripts/deploy.sh v1.0.0

# Deploy latest
./scripts/deploy.sh latest
```

## 📊 Monitoring

### Check Status
```bash
# Overall status
./scripts/monitor.sh status

# Detailed monitoring
./scripts/monitor.sh all

# Service-specific logs
./scripts/monitor.sh logs photo-svc
./scripts/monitor.sh follow user-svc

# Health checks
./scripts/monitor.sh health
```

### Resource Usage
```bash
# Resource usage
./scripts/monitor.sh resources

# Pod details
./scripts/monitor.sh pods photo-svc
```

## 🔄 Rollback

### Rollback to Previous Version
```bash
# Rollback all services
./scripts/rollback.sh

# Rollback specific service
./scripts/rollback.sh photo-svc

# Rollback to specific revision
./scripts/rollback.sh photo-svc 2
```

### View Rollout History
```bash
# All services
./scripts/rollback.sh history

# Specific service
./scripts/rollback.sh history photo-svc
```

## 🏥 Health Checks

All services include health check endpoints:

- **photo-svc**: `http://localhost:8001/health`
- **user-svc**: `http://localhost:8003/health`
- **transaction-svc**: `http://localhost:8005/health`
- **upload-svc**: `http://localhost:8002/health`
- **notification-svc**: `http://localhost:8004/health`

## 🔧 Troubleshooting

### Common Issues

1. **Deployment fails:**
   ```bash
   # Check pod status
   kubectl get pods -n be-yourmoments
   
   # Check logs
   kubectl logs deployment/photo-svc -n be-yourmoments
   ```

2. **Services not starting:**
   ```bash
   # Check events
   kubectl get events -n be-yourmoments
   
   # Check resource usage
   kubectl top pods -n be-yourmoments
   ```

3. **Database connection issues:**
   ```bash
   # Check database pod
   kubectl get pods -l app=postgres -n be-yourmoments
   
   # Check database logs
   kubectl logs deployment/postgres -n be-yourmoments
   ```

### Debugging Commands

```bash
# Follow logs for specific service
kubectl logs -f deployment/photo-svc -n be-yourmoments

# Check service endpoints
kubectl get services -n be-yourmoments

# Port forward for local testing
kubectl port-forward service/photo-svc 8001:8001 -n be-yourmoments
```

## 📈 Scaling

### Horizontal Scaling
```bash
# Scale specific service
kubectl scale deployment photo-svc --replicas=3 -n be-yourmoments

# Scale all services
for service in photo-svc user-svc transaction-svc upload-svc notification-svc; do
  kubectl scale deployment $service --replicas=3 -n be-yourmoments
done
```

### Resource Limits
Update resource limits in the Kubernetes manifests:

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "200m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

## 🔐 Security

### Secrets Management
- All sensitive data is stored in Kubernetes secrets
- Secrets are base64 encoded
- Use external secret management for production

### Network Policies
Consider implementing network policies for additional security:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

## 📝 Maintenance

### Regular Tasks
1. **Update base images** in Dockerfiles
2. **Rotate secrets** regularly
3. **Monitor resource usage**
4. **Update dependencies**

### Backup
```bash
# Backup database
kubectl exec -it deployment/postgres -n be-yourmoments -- pg_dump -U postgres postgres > backup.sql

# Backup secrets
kubectl get secrets -n be-yourmoments -o yaml > secrets-backup.yaml
```

## 📚 Documentation

- **scripts/README.md** - Detailed script documentation
- **QUICK_START.md** - Quick start guide
- **DEPLOYMENT_CHECKLIST.md** - Deployment checklist
- **TROUBLESHOOTING.md** - Troubleshooting guide

## 🆘 Support

For issues and questions:

1. Check the logs: `./scripts/monitor.sh logs [service]`
2. Check events: `./scripts/monitor.sh events`
3. Verify health: `./scripts/monitor.sh health`
4. Review deployment status: `./scripts/monitor.sh status`

## 🎯 Next Steps

1. **Configure GitHub repository secrets**
2. **Setup K3s on your VPS**
3. **Run initial setup: `./scripts/setup.sh`**
4. **Create secrets: `./scripts/create-secrets.sh`**
5. **Deploy services: `./scripts/deploy.sh latest`**
6. **Monitor: `./scripts/monitor.sh status`**

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
