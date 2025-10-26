#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

NAMESPACE="be-yourmoments"
SERVICES=("photo-svc" "user-svc" "transaction-svc" "upload-svc" "notification-svc")

echo -e "${BLUE}🔄 Rollback script for Be Your Moments microservices${NC}"

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed or not in PATH${NC}"
        exit 1
    fi
}

# Function to rollback a specific service
rollback_service() {
    local service=$1
    local revision=${2:-1}
    
    echo -e "${YELLOW}🔄 Rolling back $service to revision $revision${NC}"
    
    # Check if deployment exists
    if ! kubectl get deployment $service -n $NAMESPACE &> /dev/null; then
        echo -e "${RED}❌ Deployment $service not found${NC}"
        return 1
    fi
    
    # Rollback deployment
    kubectl rollout undo deployment/$service -n $NAMESPACE --to-revision=$revision
    
    # Wait for rollout to complete
    echo -e "${YELLOW}⏳ Waiting for $service rollout to complete...${NC}"
    kubectl rollout status deployment/$service -n $NAMESPACE --timeout=300s
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $service rolled back successfully${NC}"
    else
        echo -e "${RED}❌ $service rollback failed${NC}"
        return 1
    fi
}

# Function to rollback all services
rollback_all() {
    local revision=${1:-1}
    
    echo -e "${YELLOW}🔄 Rolling back all services to revision $revision${NC}"
    
    for service in "${SERVICES[@]}"; do
        rollback_service $service $revision
    done
    
    echo -e "${GREEN}🎉 All services rolled back successfully!${NC}"
}

# Function to show rollout history
show_history() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Rollout history for all services:${NC}"
        for service in "${SERVICES[@]}"; do
            echo -e "${BLUE}--- $service ---${NC}"
            kubectl rollout history deployment/$service -n $NAMESPACE
            echo
        done
    else
        echo -e "${YELLOW}📋 Rollout history for $service:${NC}"
        kubectl rollout history deployment/$service -n $NAMESPACE
    fi
}

# Function to show current status
show_status() {
    echo -e "${BLUE}📊 Current deployment status:${NC}"
    kubectl get deployments -n $NAMESPACE
    kubectl get pods -n $NAMESPACE
}

# Function to pause rollout
pause_rollout() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}⏸️  Pausing rollouts for all services${NC}"
        for service in "${SERVICES[@]}"; do
            kubectl rollout pause deployment/$service -n $NAMESPACE
        done
    else
        echo -e "${YELLOW}⏸️  Pausing rollout for $service${NC}"
        kubectl rollout pause deployment/$service -n $NAMESPACE
    fi
}

# Function to resume rollout
resume_rollout() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}▶️  Resuming rollouts for all services${NC}"
        for service in "${SERVICES[@]}"; do
            kubectl rollout resume deployment/$service -n $NAMESPACE
        done
    else
        echo -e "${YELLOW}▶️  Resuming rollout for $service${NC}"
        kubectl rollout resume deployment/$service -n $NAMESPACE
    fi
}

# Function to show help
show_help() {
    echo -e "${BLUE}Usage: $0 [command] [options]${NC}"
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  rollback [service] [revision]  - Rollback specific service (default: all services, revision 1)"
    echo -e "  history [service]              - Show rollout history"
    echo -e "  status                        - Show current deployment status"
    echo -e "  pause [service]               - Pause rollout"
    echo -e "  resume [service]              - Resume rollout"
    echo -e "  help                          - Show this help message"
    echo -e "${YELLOW}Examples:${NC}"
    echo -e "  $0 rollback photo-svc 2       - Rollback photo-svc to revision 2"
    echo -e "  $0 rollback                  - Rollback all services to previous revision"
    echo -e "  $0 history photo-svc          - Show history for photo-svc"
    echo -e "  $0 pause photo-svc            - Pause rollout for photo-svc"
}

# Main function
main() {
    check_kubectl
    
    case "${1:-rollback}" in
        "rollback")
            if [ -n "$2" ]; then
                rollback_service $2 $3
            else
                rollback_all $2
            fi
            show_status
            ;;
        "history")
            show_history $2
            ;;
        "status")
            show_status
            ;;
        "pause")
            pause_rollout $2
            ;;
        "resume")
            resume_rollout $2
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $1${NC}"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
