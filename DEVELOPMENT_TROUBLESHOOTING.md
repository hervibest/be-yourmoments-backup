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
