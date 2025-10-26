#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Complete setup for Be Your Moments CI/CD Pipeline${NC}"

# Function to create directory structure
create_directories() {
    echo -e "${YELLOW}📁 Creating directory structure...${NC}"
    
    mkdir -p .github/workflows
    mkdir -p scripts/k8s
    mkdir -p k8s-manifests
    
    echo -e "${GREEN}✅ Directory structure created${NC}"
}

# Function to create all necessary files
create_files() {
    echo -e "${YELLOW}📝 Creating all necessary files...${NC}"
    
    # Make all scripts executable
    chmod +x scripts/*.sh
    
    # Create .gitignore for scripts
    cat > scripts/.gitignore << 'EOF'
# Ignore sensitive files
secrets.yaml
*.key
*.pem
*.crt
*.p12
*.pfx

# Ignore temporary files
*.tmp
*.log
*.bak
EOF

    echo -e "${GREEN}✅ All files created${NC}"
}

# Function to create deployment checklist
create_checklist() {
    echo -e "${YELLOW}📋 Creating deployment checklist...${NC}"
    
    cat > DEPLOYMENT_CHECKLIST.md << 'EOF'
# Deployment Checklist

## Pre-deployment

- [ ] VPS with K3s installed and configured
- [ ] kubectl configured to access K3s cluster
- [ ] GitHub repository secrets configured
- [ ] Docker images built and pushed to registry
- [ ] Database and infrastructure services running
- [ ] Secrets created and configured

## Deployment Steps

- [ ] Run initial setup: `./scripts/setup.sh`
- [ ] Create secrets: `./scripts/create-secrets.sh`
- [ ] Deploy infrastructure: `kubectl apply -f scripts/k8s/infrastructure.yaml`
- [ ] Deploy services: `./scripts/deploy.sh latest`
- [ ] Verify deployment: `./scripts/monitor.sh status`
- [ ] Check health endpoints: `./scripts/monitor.sh health`

## Post-deployment

- [ ] All services are running
- [ ] Health checks passing
- [ ] Services can communicate with each other
- [ ] Database connections working
- [ ] External services (Redis, NATS, etc.) accessible
- [ ] Monitoring and logging working

## Rollback Plan

- [ ] Know how to rollback: `./scripts/rollback.sh`
- [ ] Have previous working version identified
- [ ] Test rollback procedure
- [ ] Document rollback steps

## Monitoring

- [ ] Set up monitoring dashboards
- [ ] Configure alerting
- [ ] Test log aggregation
- [ ] Verify metrics collection
EOF

    echo -e "${GREEN}✅ Deployment checklist created${NC}"
}

# Function to create troubleshooting guide
create_troubleshooting() {
    echo -e "${YELLOW}🔧 Creating troubleshooting guide...${NC}"
    
    cat > TROUBLESHOOTING.md << 'EOF'
# Troubleshooting Guide

## Common Issues

### 1. Services Not Starting

**Symptoms:**
- Pods stuck in Pending or CrashLoopBackOff
- Services not responding to health checks

**Solutions:**
```bash
# Check pod status
kubectl get pods -n be-yourmoments

# Check pod logs
kubectl logs deployment/photo-svc -n be-yourmoments

# Check events
kubectl get events -n be-yourmoments --sort-by='.lastTimestamp'

# Check resource usage
kubectl top pods -n be-yourmoments
```

### 2. Database Connection Issues

**Symptoms:**
- Database connection timeouts
- Authentication failures

**Solutions:**
```bash
# Check database pod
kubectl get pods -l app=postgres -n be-yourmoments

# Check database logs
kubectl logs deployment/postgres -n be-yourmoments

# Test database connection
kubectl exec -it deployment/postgres -n be-yourmoments -- psql -U postgres -d postgres

# Check secrets
kubectl get secrets -n be-yourmoments
kubectl describe secret photo-svc-secrets -n be-yourmoments
```

### 3. Service Discovery Issues

**Symptoms:**
- Services cannot find each other
- DNS resolution failures

**Solutions:**
```bash
# Check service endpoints
kubectl get endpoints -n be-yourmoments

# Check service DNS
kubectl exec -it deployment/photo-svc -n be-yourmoments -- nslookup user-svc

# Check Consul
kubectl logs deployment/consul -n be-yourmoments
```

### 4. Resource Issues

**Symptoms:**
- Pods evicted due to resource constraints
- High CPU/memory usage

**Solutions:**
```bash
# Check resource usage
kubectl top pods -n be-yourmoments
kubectl top nodes

# Check resource limits
kubectl describe pod <pod-name> -n be-yourmoments

# Scale services
kubectl scale deployment photo-svc --replicas=2 -n be-yourmoments
```

### 5. Network Issues

**Symptoms:**
- Services cannot communicate
- Port forwarding not working

**Solutions:**
```bash
# Check services
kubectl get services -n be-yourmoments

# Test port forwarding
kubectl port-forward service/photo-svc 8001:8001 -n be-yourmoments

# Check network policies
kubectl get networkpolicies -n be-yourmoments
```

## Debugging Commands

### General Debugging
```bash
# Check overall status
./scripts/monitor.sh status

# Check specific service
./scripts/monitor.sh logs photo-svc

# Follow logs
./scripts/monitor.sh follow user-svc

# Check health
./scripts/monitor.sh health
```

### Kubernetes Debugging
```bash
# Describe resources
kubectl describe pod <pod-name> -n be-yourmoments
kubectl describe service <service-name> -n be-yourmoments
kubectl describe deployment <deployment-name> -n be-yourmoments

# Check logs
kubectl logs <pod-name> -n be-yourmoments
kubectl logs <pod-name> -n be-yourmoments --previous

# Check events
kubectl get events -n be-yourmoments --sort-by='.lastTimestamp'
```

### Application Debugging
```bash
# Check application logs
kubectl logs deployment/photo-svc -n be-yourmoments

# Check configuration
kubectl exec -it deployment/photo-svc -n be-yourmoments -- env

# Test connectivity
kubectl exec -it deployment/photo-svc -n be-yourmoments -- curl http://user-svc:8003/health
```

## Performance Issues

### High CPU Usage
```bash
# Check CPU usage
kubectl top pods -n be-yourmoments

# Check CPU limits
kubectl describe pod <pod-name> -n be-yourmoments | grep -A 5 "Limits:"

# Scale horizontally
kubectl scale deployment photo-svc --replicas=3 -n be-yourmoments
```

### High Memory Usage
```bash
# Check memory usage
kubectl top pods -n be-yourmoments

# Check memory limits
kubectl describe pod <pod-name> -n be-yourmoments | grep -A 5 "Limits:"

# Check for memory leaks
kubectl logs deployment/photo-svc -n be-yourmoments | grep -i "memory\|oom"
```

### Slow Response Times
```bash
# Check service response times
kubectl exec -it deployment/photo-svc -n be-yourmoments -- curl -w "@curl-format.txt" -o /dev/null -s http://user-svc:8003/health

# Check database performance
kubectl exec -it deployment/postgres -n be-yourmoments -- psql -U postgres -c "SELECT * FROM pg_stat_activity;"
```

## Recovery Procedures

### Service Recovery
```bash
# Restart service
kubectl rollout restart deployment/photo-svc -n be-yourmoments

# Check rollout status
kubectl rollout status deployment/photo-svc -n be-yourmoments
```

### Database Recovery
```bash
# Restart database
kubectl rollout restart deployment/postgres -n be-yourmoments

# Check database status
kubectl exec -it deployment/postgres -n be-yourmoments -- pg_isready
```

### Full System Recovery
```bash
# Restart all services
for service in photo-svc user-svc transaction-svc upload-svc notification-svc; do
  kubectl rollout restart deployment/$service -n be-yourmoments
done

# Check all services
./scripts/monitor.sh status
```

## Contact Information

For additional support:
- Check the logs first
- Review the monitoring dashboard
- Consult the documentation
- Contact the development team
EOF

    echo -e "${GREEN}✅ Troubleshooting guide created${NC}"
}

# Function to show final summary
show_summary() {
    echo -e "${GREEN}🎉 Complete setup finished!${NC}"
    echo -e "${BLUE}📋 Summary of created files:${NC}"
    echo -e "${YELLOW}CI/CD Pipeline:${NC}"
    echo -e "  - .github/workflows/ci-cd.yml"
    echo -e "  - scripts/deploy.sh"
    echo -e "  - scripts/setup.sh"
    echo -e "  - scripts/rollback.sh"
    echo -e "  - scripts/monitor.sh"
    echo -e "  - scripts/create-secrets.sh"
    echo -e "  - scripts/setup-env.sh"
    echo -e "  - scripts/k3s-setup.sh"
    echo
    echo -e "${YELLOW}Kubernetes Manifests:${NC}"
    echo -e "  - scripts/k8s/photo-svc.yaml"
    echo -e "  - scripts/k8s/user-svc.yaml"
    echo -e "  - scripts/k8s/transaction-svc.yaml"
    echo -e "  - scripts/k8s/upload-svc.yaml"
    echo -e "  - scripts/k8s/notification-svc.yaml"
    echo -e "  - scripts/k8s/infrastructure.yaml"
    echo -e "  - scripts/k8s/secrets.yaml"
    echo -e "  - scripts/k8s/ingress.yaml"
    echo
    echo -e "${YELLOW}Documentation:${NC}"
    echo -e "  - scripts/README.md"
    echo -e "  - QUICK_START.md"
    echo -e "  - DEPLOYMENT_CHECKLIST.md"
    echo -e "  - TROUBLESHOOTING.md"
    echo
    echo -e "${BLUE}🚀 Next steps:${NC}"
    echo -e "1. Configure GitHub repository secrets"
    echo -e "2. Setup K3s on your VPS"
    echo -e "3. Run: ./scripts/setup.sh"
    echo -e "4. Create secrets: ./scripts/create-secrets.sh"
    echo -e "5. Deploy: ./scripts/deploy.sh latest"
    echo -e "6. Monitor: ./scripts/monitor.sh status"
}

# Main function
main() {
    create_directories
    create_files
    create_checklist
    create_troubleshooting
    show_summary
}

# Run main function
main "$@"
