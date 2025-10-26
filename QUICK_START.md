# Quick Start Guide

## Prerequisites

1. **VPS with K3s installed**
   ```bash
   # On your VPS, run:
   curl -sfL https://get.k3s.io | sh -
   ```

2. **GitHub repository with secrets configured**
   - Go to Settings > Secrets and variables > Actions
   - Add: `SSH_PRIVATE_KEY`, `VPS_USER`, `VPS_HOST`

## Setup Steps

### 1. Initial Setup

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Setup infrastructure
./scripts/setup.sh

# Create secrets
./scripts/create-secrets.sh interactive
```

### 2. Deploy Services

```bash
# Deploy all services
./scripts/deploy.sh latest

# Check status
./scripts/monitor.sh status
```

### 3. Create and Push Tag

```bash
# Create a tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically build and deploy
```

## Monitoring

```bash
# Check overall status
./scripts/monitor.sh status

# Check specific service logs
./scripts/monitor.sh logs photo-svc

# Follow logs in real-time
./scripts/monitor.sh follow user-svc

# Check health endpoints
./scripts/monitor.sh health
```

## Rollback

```bash
# Rollback all services
./scripts/rollback.sh

# Rollback specific service
./scripts/rollback.sh photo-svc

# View rollout history
./scripts/rollback.sh history
```

## Troubleshooting

```bash
# Check pod status
kubectl get pods -n be-yourmoments

# Check events
kubectl get events -n be-yourmoments

# Check logs
kubectl logs deployment/photo-svc -n be-yourmoments
```
