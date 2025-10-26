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
