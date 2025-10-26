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
