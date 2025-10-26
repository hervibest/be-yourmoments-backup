# Be Your Moments CI/CD Pipeline

This repository contains a complete CI/CD pipeline for the Be Your Moments microservices application, designed to work with GitHub Actions and K3s.

## 🏗️ Architecture

The application consists of 5 microservices:
- **photo-svc** - Photo management service
- **user-svc** - User management service  
- **transaction-svc** - Transaction processing service
- **upload-svc** - File upload service
- **notification-svc** - Notification service

## 🚀 CI/CD Pipeline

### GitHub Actions Workflow

The CI/CD pipeline is triggered by:
- **Push to tags** (e.g., `v1.0.0`) - Builds and deploys to production
- **Pull requests** - Runs tests and linting

### Pipeline Stages

1. **Test Stage**
   - Runs unit tests for all services
   - Performs linting with golangci-lint
   - Caches Go modules for faster builds

2. **Build Stage** (on tag push)
   - Builds Docker images for all services
   - Pushes images to GitHub Container Registry
   - Uses multi-stage builds for optimized images

3. **Deploy Stage** (on tag push)
   - SSH to VPS
   - Updates Kubernetes deployments
   - Performs rolling updates

## 📁 Scripts Overview

### Core Scripts

- **`deploy.sh`** - Main deployment script
- **`setup.sh`** - Initial infrastructure setup
- **`rollback.sh`** - Rollback deployments
- **`monitor.sh`** - Monitoring and troubleshooting

### Kubernetes Manifests

- **`k8s/*.yaml`** - Service-specific Kubernetes manifests
- **`k8s/infrastructure.yaml`** - Infrastructure components (PostgreSQL, Redis, NATS, Consul, MinIO)
- **`k8s/secrets.yaml`** - Secrets configuration

## 🛠️ Setup Instructions

### 1. Prerequisites

- K3s cluster running on your VPS
- kubectl configured to access your cluster
- SSH access to your VPS
- GitHub repository with secrets configured

### 2. GitHub Secrets

Configure the following secrets in your GitHub repository:

```
SSH_PRIVATE_KEY    # Private key for SSH access to VPS
VPS_USER          # Username for VPS SSH access
VPS_HOST          # IP address or hostname of your VPS
```

### 3. Initial Setup

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Setup infrastructure
./scripts/setup.sh

# Update secrets with your actual values
# Edit scripts/k8s/secrets.yaml

# Deploy services
./scripts/deploy.sh v1.0.0
```

### 4. Update Secrets

Before deploying, update the secrets in `scripts/k8s/secrets.yaml`:

```bash
# Encode your values in base64
echo -n "your-db-password" | base64
echo -n "your-redis-password" | base64

# Update the secrets.yaml file with your encoded values
```

## 🚀 Deployment Process

### Automatic Deployment

1. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will:**
   - Run tests
   - Build Docker images
   - Push to registry
   - Deploy to K3s

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

### Logs and Debugging

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

## 🆘 Support

For issues and questions:

1. Check the logs: `./scripts/monitor.sh logs [service]`
2. Check events: `./scripts/monitor.sh events`
3. Verify health: `./scripts/monitor.sh health`
4. Review deployment status: `./scripts/monitor.sh status`
