#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Complete Development Environment Setup${NC}"

# Function to create all necessary files
create_files() {
    echo -e "${YELLOW}📝 Creating all necessary files...${NC}"
    
    # Make all scripts executable
    chmod +x scripts/*.sh
    
    # Create .gitignore for development
    cat > .gitignore-dev << 'EOF'
# Development environment
/tmp/be-yourmoments-dev/
*.log
*.tmp

# Service account keys
**/serviceAccountKey.json
**/service-account*.json

# Environment files
.env
.env.*
!.env.example

# Docker volumes
docker-data/
postgres-data/
minio-data/
EOF

    echo -e "${GREEN}✅ All files created${NC}"
}

# Function to create development checklist
create_checklist() {
    echo -e "${YELLOW}📋 Creating development checklist...${NC}"
    
    cat > DEVELOPMENT_CHECKLIST.md << 'EOF'
# Development Environment Checklist

## Pre-development

- [ ] Docker and Docker Compose installed
- [ ] Git repository cloned
- [ ] Service account keys available
- [ ] Development environment setup

## Development Setup

- [ ] Run setup script: `./scripts/setup-dev.sh`
- [ ] Update service account keys
- [ ] Start services: `cd /tmp/be-yourmoments-dev && ./start.sh`
- [ ] Verify all services are running
- [ ] Check health endpoints

## Development Workflow

- [ ] Make code changes
- [ ] Test locally
- [ ] Create development tag: `git tag x.dev.1.0.0`
- [ ] Push tag: `git push origin x.dev.1.0.0`
- [ ] Monitor deployment
- [ ] Verify services on VPS

## Monitoring

- [ ] Check service status
- [ ] Monitor resource usage
- [ ] Check logs for errors
- [ ] Verify health endpoints
- [ ] Test service communication

## Cleanup

- [ ] Stop services when done
- [ ] Clean up unused resources
- [ ] Remove old images
- [ ] Clean up volumes if needed
EOF

    echo -e "${GREEN}✅ Development checklist created${NC}"
}

# Function to create troubleshooting guide
create_troubleshooting() {
    echo -e "${YELLOW}🔧 Creating troubleshooting guide...${NC}"
    
    cat > DEVELOPMENT_TROUBLESHOOTING.md << 'EOF'
# Development Environment Troubleshooting

## Common Issues

### 1. Services Not Starting

**Symptoms:**
- Containers fail to start
- Services not responding
- Port conflicts

**Solutions:**
```bash
# Check Docker status
docker ps -a

# Check logs
./logs.sh

# Restart services
./restart.sh

# Check port conflicts
netstat -tulpn | grep -E "(8001|8002|8003|8004|8005)"
```

### 2. Database Connection Issues

**Symptoms:**
- Database connection timeouts
- Migration failures
- Service startup failures

**Solutions:**
```bash
# Check PostgreSQL status
docker exec -it be-yourmoments-dev_postgres_1 pg_isready

# Check database logs
docker logs be-yourmoments-dev_postgres_1

# Restart database
docker restart be-yourmoments-dev_postgres_1
```

### 3. Service Discovery Issues

**Symptoms:**
- Services cannot find each other
- Consul registration failures

**Solutions:**
```bash
# Check Consul status
curl http://localhost:8500/v1/status/leader

# Check service registration
curl http://localhost:8500/v1/catalog/services

# Restart Consul
docker restart be-yourmoments-dev_consul_1
```

### 4. Resource Issues

**Symptoms:**
- High memory usage
- Slow performance
- Container evictions

**Solutions:**
```bash
# Check resource usage
./scripts/monitor-dev.sh resources

# Check Docker stats
docker stats

# Clean up unused resources
docker system prune -f
```

## Debugging Commands

### General Debugging
```bash
# Check overall status
./scripts/monitor-dev.sh status

# Check specific service
./scripts/monitor-dev.sh logs photo-svc

# Follow logs
./scripts/monitor-dev.sh follow user-svc

# Check health
./scripts/monitor-dev.sh health
```

### Container Debugging
```bash
# Show all containers
docker ps -a

# Show container logs
docker logs <container_name>

# Execute command in container
docker exec -it <container_name> bash

# Show container details
docker inspect <container_name>
```

### Network Debugging
```bash
# Show networks
docker network ls

# Show network details
docker network inspect backend

# Test connectivity
docker exec -it <container_name> ping <other_container>
```

### Storage Debugging
```bash
# Show volumes
docker volume ls

# Show volume details
docker volume inspect <volume_name>

# Check disk usage
docker system df
```

## Performance Issues

### High Memory Usage
```bash
# Check memory usage
./scripts/monitor-dev.sh resources

# Check specific container
docker stats <container_name>

# Restart high-memory containers
docker restart <container_name>
```

### Slow Response Times
```bash
# Check service response times
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8001/health

# Check database performance
docker exec -it be-yourmoments-dev_postgres_1 psql -U postgres -c "SELECT * FROM pg_stat_activity;"
```

### CPU Issues
```bash
# Check CPU usage
docker stats --no-stream

# Check specific container CPU
docker exec -it <container_name> top
```

## Recovery Procedures

### Service Recovery
```bash
# Restart specific service
./scripts/monitor-dev.sh restart photo-svc

# Restart all services
./restart.sh

# Check service status
./status.sh
```

### Database Recovery
```bash
# Restart database
docker restart be-yourmoments-dev_postgres_1

# Check database status
docker exec -it be-yourmoments-dev_postgres_1 pg_isready

# Run migrations
docker exec -it be-yourmoments-dev_photo-svc_1 ./migrate up
```

### Full System Recovery
```bash
# Stop all services
./stop.sh

# Clean up resources
docker system prune -f

# Restart services
./start.sh

# Check status
./status.sh
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
    echo -e "${GREEN}🎉 Development environment setup completed!${NC}"
    echo -e "${BLUE}📋 Summary of created files:${NC}"
    echo -e "${YELLOW}CI/CD Pipeline:${NC}"
    echo -e "  - .github/workflows/ci-cd-dev.yml"
    echo -e "  - scripts/deploy-dev.sh"
    echo -e "  - scripts/setup-dev.sh"
    echo -e "  - scripts/monitor-dev.sh"
    echo -e "  - scripts/complete-dev-setup.sh"
    echo
    echo -e "${YELLOW}Documentation:${NC}"
    echo -e "  - DEVELOPMENT_README.md"
    echo -e "  - DEVELOPMENT_CHECKLIST.md"
    echo -e "  - DEVELOPMENT_TROUBLESHOOTING.md"
    echo
    echo -e "${BLUE}🚀 Next steps:${NC}"
    echo -e "1. Setup development environment: ./scripts/setup-dev.sh"
    echo -e "2. Start services: cd /tmp/be-yourmoments-dev && ./start.sh"
    echo -e "3. Monitor services: ./scripts/monitor-dev.sh status"
    echo -e "4. Create development tag: git tag x.dev.1.0.0"
    echo -e "5. Deploy to VPS: git push origin x.dev.1.0.0"
    echo
    echo -e "${YELLOW}🌐 Service URLs:${NC}"
    echo -e "Photo Service: http://localhost:8001"
    echo -e "User Service: http://localhost:8003"
    echo -e "Transaction Service: http://localhost:8005"
    echo -e "Upload Service: http://localhost:8002"
    echo -e "Notification Service: http://localhost:8004"
    echo -e "MinIO Console: http://localhost:9001"
    echo -e "Consul UI: http://localhost:8500"
}

# Main function
main() {
    create_files
    create_checklist
    create_troubleshooting
    show_summary
}

# Run main function
main "$@"
