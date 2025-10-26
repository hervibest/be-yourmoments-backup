#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REGISTRY="ghcr.io"
REPOSITORY="hervipro/be-yourmoments-backup"
TAG=${1:-latest}
NAMESPACE="be-yourmoments"

# Services
SERVICES=("photo-svc" "user-svc" "transaction-svc" "upload-svc" "notification-svc")

echo -e "${GREEN}🚀 Starting deployment for tag: $TAG${NC}"

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed or not in PATH${NC}"
        exit 1
    fi
}

# Function to check if namespace exists
check_namespace() {
    if ! kubectl get namespace $NAMESPACE &> /dev/null; then
        echo -e "${YELLOW}📦 Creating namespace: $NAMESPACE${NC}"
        kubectl create namespace $NAMESPACE
    else
        echo -e "${GREEN}✅ Namespace $NAMESPACE already exists${NC}"
    fi
}

# Function to deploy a service
deploy_service() {
    local service=$1
    local image="$REGISTRY/$REPOSITORY-$service:$TAG"
    
    echo -e "${YELLOW}🔄 Deploying $service with image: $image${NC}"
    
    # Update image in deployment
    kubectl set image deployment/$service $service=$image -n $NAMESPACE
    
    # Wait for rollout to complete
    echo -e "${YELLOW}⏳ Waiting for $service rollout to complete...${NC}"
    kubectl rollout status deployment/$service -n $NAMESPACE --timeout=300s
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $service deployed successfully${NC}"
    else
        echo -e "${RED}❌ $service deployment failed${NC}"
        exit 1
    fi
}

# Function to check if deployment exists
check_deployment() {
    local service=$1
    if kubectl get deployment $service -n $NAMESPACE &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Function to create deployment if it doesn't exist
create_deployment() {
    local service=$1
    local image="$REGISTRY/$REPOSITORY-$service:$TAG"
    
    echo -e "${YELLOW}📦 Creating deployment for $service${NC}"
    
    # Apply the Kubernetes manifest
    kubectl apply -f ${service}.yaml -n $NAMESPACE
    
    # Wait for deployment to be ready
    kubectl wait --for=condition=available --timeout=300s deployment/$service -n $NAMESPACE
}

# Main deployment logic
main() {
    echo -e "${GREEN}🎯 Starting deployment process${NC}"
    
    # Check prerequisites
    check_kubectl
    check_namespace
    
    # Deploy each service
    for service in "${SERVICES[@]}"; do
        if check_deployment $service; then
            deploy_service $service
        else
            create_deployment $service
        fi
    done
    
    echo -e "${GREEN}🎉 All services deployed successfully!${NC}"
    
    # Show status
    echo -e "${YELLOW}📊 Deployment status:${NC}"
    kubectl get deployments -n $NAMESPACE
    kubectl get pods -n $NAMESPACE
    kubectl get services -n $NAMESPACE
}

# Run main function
main "$@"
