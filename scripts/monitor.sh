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

echo -e "${BLUE}📊 Monitoring script for Be Your Moments microservices${NC}"

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed or not in PATH${NC}"
        exit 1
    fi
}

# Function to show overall status
show_status() {
    echo -e "${BLUE}📊 Overall Status${NC}"
    echo -e "${YELLOW}Namespaces:${NC}"
    kubectl get namespaces | grep $NAMESPACE || echo "Namespace not found"
    
    echo -e "${YELLOW}Deployments:${NC}"
    kubectl get deployments -n $NAMESPACE
    
    echo -e "${YELLOW}Services:${NC}"
    kubectl get services -n $NAMESPACE
    
    echo -e "${YELLOW}Pods:${NC}"
    kubectl get pods -n $NAMESPACE
}

# Function to show detailed pod status
show_pod_details() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Pod details for all services:${NC}"
        for service in "${SERVICES[@]}"; do
            echo -e "${BLUE}--- $service ---${NC}"
            kubectl get pods -l app=$service -n $NAMESPACE -o wide
            echo
        done
    else
        echo -e "${YELLOW}📋 Pod details for $service:${NC}"
        kubectl get pods -l app=$service -n $NAMESPACE -o wide
    fi
}

# Function to show resource usage
show_resources() {
    echo -e "${YELLOW}💾 Resource usage:${NC}"
    kubectl top pods -n $NAMESPACE
    echo
    kubectl top nodes
}

# Function to show logs
show_logs() {
    local service=$1
    local lines=${2:-50}
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Showing last $lines lines of logs for $service:${NC}"
    kubectl logs deployment/$service -n $NAMESPACE --tail=$lines
}

# Function to follow logs
follow_logs() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Following logs for $service (Ctrl+C to stop):${NC}"
    kubectl logs -f deployment/$service -n $NAMESPACE
}

# Function to show events
show_events() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Recent events in namespace $NAMESPACE:${NC}"
        kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp' | tail -20
    else
        echo -e "${YELLOW}📋 Recent events for $service:${NC}"
        kubectl get events -n $NAMESPACE --field-selector involvedObject.name=$service --sort-by='.lastTimestamp'
    fi
}

# Function to check health endpoints
check_health() {
    echo -e "${YELLOW}🏥 Health check for all services:${NC}"
    
    for service in "${SERVICES[@]}"; do
        echo -e "${BLUE}Checking $service...${NC}"
        
        # Get service port
        local port=$(kubectl get service $service -n $NAMESPACE -o jsonpath='{.spec.ports[0].port}' 2>/dev/null || echo "N/A")
        
        if [ "$port" != "N/A" ]; then
            # Port forward and check health
            kubectl port-forward service/$service $port:$port -n $NAMESPACE > /dev/null 2>&1 &
            local pf_pid=$!
            sleep 2
            
            if curl -s http://localhost:$port/health > /dev/null 2>&1; then
                echo -e "${GREEN}✅ $service is healthy${NC}"
            else
                echo -e "${RED}❌ $service is unhealthy${NC}"
            fi
            
            kill $pf_pid 2>/dev/null || true
        else
            echo -e "${YELLOW}⚠️  $service service not found${NC}"
        fi
    done
}

# Function to show network policies
show_network() {
    echo -e "${YELLOW}🌐 Network policies:${NC}"
    kubectl get networkpolicies -n $NAMESPACE
    
    echo -e "${YELLOW}Ingress:${NC}"
    kubectl get ingress -n $NAMESPACE
}

# Function to show persistent volumes
show_storage() {
    echo -e "${YELLOW}💾 Storage:${NC}"
    kubectl get pv,pvc -n $NAMESPACE
}

# Function to show secrets
show_secrets() {
    echo -e "${YELLOW}🔐 Secrets:${NC}"
    kubectl get secrets -n $NAMESPACE
}

# Function to show configmaps
show_configmaps() {
    echo -e "${YELLOW}⚙️  ConfigMaps:${NC}"
    kubectl get configmaps -n $NAMESPACE
}

# Function to show help
show_help() {
    echo -e "${BLUE}Usage: $0 [command] [options]${NC}"
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  status                    - Show overall status"
    echo -e "  pods [service]            - Show pod details"
    echo -e "  resources                - Show resource usage"
    echo -e "  logs [service] [lines]    - Show logs (default: 50 lines)"
    echo -e "  follow [service]          - Follow logs in real-time"
    echo -e "  events [service]          - Show recent events"
    echo -e "  health                    - Check health endpoints"
    echo -e "  network                   - Show network policies and ingress"
    echo -e "  storage                   - Show persistent volumes"
    echo -e "  secrets                   - Show secrets"
    echo -e "  configmaps                - Show configmaps"
    echo -e "  all                       - Show comprehensive status"
    echo -e "  help                      - Show this help message"
    echo -e "${YELLOW}Examples:${NC}"
    echo -e "  $0 status                - Show overall status"
    echo -e "  $0 logs photo-svc 100    - Show last 100 lines of photo-svc logs"
    echo -e "  $0 follow user-svc        - Follow user-svc logs"
    echo -e "  $0 health                - Check all service health"
}

# Function to show comprehensive status
show_all() {
    show_status
    echo
    show_resources
    echo
    show_network
    echo
    show_storage
    echo
    show_secrets
    echo
    show_configmaps
    echo
    show_events
}

# Main function
main() {
    check_kubectl
    
    case "${1:-status}" in
        "status")
            show_status
            ;;
        "pods")
            show_pod_details $2
            ;;
        "resources")
            show_resources
            ;;
        "logs")
            show_logs $2 $3
            ;;
        "follow")
            follow_logs $2
            ;;
        "events")
            show_events $2
            ;;
        "health")
            check_health
            ;;
        "network")
            show_network
            ;;
        "storage")
            show_storage
            ;;
        "secrets")
            show_secrets
            ;;
        "configmaps")
            show_configmaps
            ;;
        "all")
            show_all
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
