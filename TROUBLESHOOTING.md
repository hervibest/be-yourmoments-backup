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
