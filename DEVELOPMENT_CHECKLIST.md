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
