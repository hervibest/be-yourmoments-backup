#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Setting up environment for Be Your Moments CI/CD${NC}"

# Function to make scripts executable
make_executable() {
    echo -e "${YELLOW}🔧 Making scripts executable...${NC}"
    chmod +x scripts/*.sh
    echo -e "${GREEN}✅ Scripts are now executable${NC}"
}

# Function to create .env template
create_env_template() {
    echo -e "${YELLOW}📝 Creating .env template...${NC}"
    
    cat > .env.template << 'EOF'
# Database Configuration
DB_HOST=postgres-service
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# Redis Configuration
REDIS_HOST=redis-service
REDIS_PASSWORD=your_redis_password

# NATS Configuration
NATS_URL=nats://nats-service:4222

# Consul Configuration
CONSUL_HOST=consul-service:8500

# MinIO Configuration (for upload-svc)
MINIO_ENDPOINT=http://minio-service:9000
MINIO_ACCESS_KEY=your_minio_access_key
MINIO_SECRET_KEY=your_minio_secret_key

# GitHub Container Registry
REGISTRY=ghcr.io
REPOSITORY=your-username/be-yourmoments-backup

# VPS Configuration
VPS_USER=your_vps_user
VPS_HOST=your_vps_ip
EOF

    echo -e "${GREEN}✅ .env.template created${NC}"
}

# Function to create GitHub Actions environment file
create_github_env() {
    echo -e "${YELLOW}📝 Creating GitHub Actions environment file...${NC}"
    
    cat > .github/workflows/env.example << 'EOF'
# Copy this file to your GitHub repository secrets
# Go to Settings > Secrets and variables > Actions > Repository secrets

# SSH Configuration
SSH_PRIVATE_KEY=your_ssh_private_key_here
VPS_USER=your_vps_username
VPS_HOST=your_vps_ip_address

# Database Configuration (if needed for testing)
DB_HOST=postgres-service
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# Redis Configuration
REDIS_HOST=redis-service
REDIS_PASSWORD=your_redis_password

# NATS Configuration
NATS_URL=nats://nats-service:4222

# Consul Configuration
CONSUL_HOST=consul-service:8500

# MinIO Configuration
MINIO_ENDPOINT=http://minio-service:9000
MINIO_ACCESS_KEY=your_minio_access_key
MINIO_SECRET_KEY=your_minio_secret_key
EOF

    echo -e "${GREEN}✅ GitHub Actions environment file created${NC}"
}

# Function to create K3s setup script
create_k3s_setup() {
    echo -e "${YELLOW}📝 Creating K3s setup script...${NC}"
    
    cat > scripts/k3s-setup.sh << 'EOF'
#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Setting up K3s for Be Your Moments${NC}"

# Function to install K3s
install_k3s() {
    echo -e "${YELLOW}📦 Installing K3s...${NC}"
    curl -sfL https://get.k3s.io | sh -
    
    # Add kubectl to PATH
    echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
    source ~/.bashrc
    
    # Copy kubeconfig
    mkdir -p ~/.kube
    sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
    sudo chown $(id -u):$(id -g) ~/.kube/config
    
    echo -e "${GREEN}✅ K3s installed successfully${NC}"
}

# Function to install required tools
install_tools() {
    echo -e "${YELLOW}🔧 Installing required tools...${NC}"
    
    # Install kubectl
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl
    
    # Install helm
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    
    echo -e "${GREEN}✅ Tools installed successfully${NC}"
}

# Function to setup storage
setup_storage() {
    echo -e "${YELLOW}💾 Setting up storage...${NC}"
    
    # Create local storage class
    cat << 'STORAGE_YAML' | kubectl apply -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: rancher.io/local-path
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
STORAGE_YAML
    
    echo -e "${GREEN}✅ Storage configured${NC}"
}

# Function to setup ingress
setup_ingress() {
    echo -e "${YELLOW}🌐 Setting up ingress...${NC}"
    
    # Install traefik (comes with K3s)
    kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml
    kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v2.10/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1alpha1.yml
    
    echo -e "${GREEN}✅ Ingress configured${NC}"
}

# Main setup function
main() {
    install_k3s
    install_tools
    setup_storage
    setup_ingress
    
    echo -e "${GREEN}🎉 K3s setup completed successfully!${NC}"
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "1. Configure kubectl: export KUBECONFIG=~/.kube/config"
    echo -e "2. Test connection: kubectl get nodes"
    echo -e "3. Run: ./scripts/setup.sh"
}

# Run main function
main "$@"
EOF

    chmod +x scripts/k3s-setup.sh
    echo -e "${GREEN}✅ K3s setup script created${NC}"
}

# Function to create quick start guide
create_quick_start() {
    echo -e "${YELLOW}📝 Creating quick start guide...${NC}"
    
    cat > QUICK_START.md << 'EOF'
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
EOF

    echo -e "${GREEN}✅ Quick start guide created${NC}"
}

# Function to show next steps
show_next_steps() {
    echo -e "${GREEN}🎉 Environment setup completed!${NC}"
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "1. Update .env.template with your actual values"
    echo -e "2. Configure GitHub repository secrets"
    echo -e "3. Setup K3s on your VPS: ./scripts/k3s-setup.sh"
    echo -e "4. Run initial setup: ./scripts/setup.sh"
    echo -e "5. Create secrets: ./scripts/create-secrets.sh"
    echo -e "6. Deploy services: ./scripts/deploy.sh latest"
    echo -e "7. Monitor: ./scripts/monitor.sh status"
}

# Main function
main() {
    make_executable
    create_env_template
    create_github_env
    create_k3s_setup
    create_quick_start
    show_next_steps
}

# Run main function
main "$@"
