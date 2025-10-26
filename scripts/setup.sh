#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

NAMESPACE="be-yourmoments"

echo -e "${BLUE}🚀 Setting up Be Your Moments microservices on K3s${NC}"

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed or not in PATH${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ kubectl is available${NC}"
}

# Function to create namespace
create_namespace() {
    if ! kubectl get namespace $NAMESPACE &> /dev/null; then
        echo -e "${YELLOW}📦 Creating namespace: $NAMESPACE${NC}"
        kubectl create namespace $NAMESPACE
    else
        echo -e "${GREEN}✅ Namespace $NAMESPACE already exists${NC}"
    fi
}

# Function to deploy infrastructure
deploy_infrastructure() {
    echo -e "${YELLOW}🏗️  Deploying infrastructure components...${NC}"
    
    # Deploy infrastructure
    kubectl apply -f scripts/k8s/infrastructure.yaml -n $NAMESPACE
    
    # Wait for infrastructure to be ready
    echo -e "${YELLOW}⏳ Waiting for infrastructure to be ready...${NC}"
    kubectl wait --for=condition=available --timeout=300s deployment/postgres -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/redis -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/nats -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/consul -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/minio -n $NAMESPACE
    
    echo -e "${GREEN}✅ Infrastructure deployed successfully${NC}"
}

# Function to setup secrets
setup_secrets() {
    echo -e "${YELLOW}🔐 Setting up secrets...${NC}"
    
    # Check if secrets file exists
    if [ ! -f "scripts/k8s/secrets.yaml" ]; then
        echo -e "${RED}❌ secrets.yaml not found. Please create it first.${NC}"
        exit 1
    fi
    
    # Apply secrets
    kubectl apply -f scripts/k8s/secrets.yaml -n $NAMESPACE
    
    echo -e "${GREEN}✅ Secrets configured successfully${NC}"
}

# Function to show status
show_status() {
    echo -e "${BLUE}📊 Current deployment status:${NC}"
    echo -e "${YELLOW}Namespaces:${NC}"
    kubectl get namespaces | grep $NAMESPACE || echo "Namespace not found"
    
    echo -e "${YELLOW}Deployments:${NC}"
    kubectl get deployments -n $NAMESPACE
    
    echo -e "${YELLOW}Services:${NC}"
    kubectl get services -n $NAMESPACE
    
    echo -e "${YELLOW}Pods:${NC}"
    kubectl get pods -n $NAMESPACE
}

# Function to show logs
show_logs() {
    local service=$1
    if [ -z "$service" ]; then
        echo -e "${YELLOW}Available services: postgres, redis, nats, consul, minio${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Showing logs for $service:${NC}"
    kubectl logs -f deployment/$service -n $NAMESPACE
}

# Function to cleanup
cleanup() {
    echo -e "${YELLOW}🧹 Cleaning up deployment...${NC}"
    kubectl delete namespace $NAMESPACE
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Main function
main() {
    case "${1:-setup}" in
        "setup")
            check_kubectl
            create_namespace
            deploy_infrastructure
            setup_secrets
            show_status
            echo -e "${GREEN}🎉 Setup completed successfully!${NC}"
            echo -e "${BLUE}Next steps:${NC}"
            echo -e "1. Update secrets.yaml with your actual values"
            echo -e "2. Run: ./scripts/deploy.sh to deploy services"
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs $2
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"-h"|"--help")
            echo -e "${BLUE}Usage: $0 [command]${NC}"
            echo -e "${YELLOW}Commands:${NC}"
            echo -e "  setup    - Setup infrastructure and secrets (default)"
            echo -e "  status   - Show deployment status"
            echo -e "  logs     - Show logs for a service"
            echo -e "  cleanup  - Remove all resources"
            echo -e "  help     - Show this help message"
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $1${NC}"
            echo -e "Run '$0 help' for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
