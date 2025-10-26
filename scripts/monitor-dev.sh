#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

COMPOSE_FILE="docker-compose-dev.yaml"
SERVICES=("photo-svc" "user-svc" "transaction-svc" "upload-svc" "notification-svc")

echo -e "${BLUE}📊 Monitoring Be Your Moments Development Environment${NC}"

# Function to check if Docker is running
check_docker() {
    if ! docker info &> /dev/null; then
        echo -e "${RED}❌ Docker is not running${NC}"
        exit 1
    fi
}

# Function to check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not available${NC}"
        exit 1
    fi
}

# Function to get compose command
get_compose_cmd() {
    if command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    else
        echo "docker compose"
    fi
}

# Function to show overall status
show_status() {
    echo -e "${BLUE}📊 Overall Status${NC}"
    
    cd /tmp/be-yourmoments-dev
    local compose_cmd=$(get_compose_cmd)
    
    echo -e "${YELLOW}Services:${NC}"
    $compose_cmd -f $COMPOSE_FILE ps
    
    echo -e "${YELLOW}Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
}

# Function to show detailed pod status
show_service_details() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Service details for all services:${NC}"
        for service in "${SERVICES[@]}"; do
            echo -e "${BLUE}--- $service ---${NC}"
            docker ps --filter "name=$service" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
            echo
        done
    else
        echo -e "${YELLOW}📋 Service details for $service:${NC}"
        docker ps --filter "name=$service" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    fi
}

# Function to show resource usage
show_resources() {
    echo -e "${YELLOW}💾 Resource usage:${NC}"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
    echo
    echo -e "${YELLOW}💽 Disk usage:${NC}"
    docker system df
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
    docker logs --tail=$lines $service
}

# Function to follow logs
follow_logs() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Following logs for $service (Ctrl+C to stop):${NC}"
    docker logs -f $service
}

# Function to show events
show_events() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Recent Docker events:${NC}"
        docker events --since=1h --format "table {{.Time}}\t{{.Type}}\t{{.Actor.Attributes.name}}"
    else
        echo -e "${YELLOW}📋 Recent events for $service:${NC}"
        docker events --since=1h --filter container=$service --format "table {{.Time}}\t{{.Type}}\t{{.Actor.Attributes.name}}"
    fi
}

# Function to check health endpoints
check_health() {
    echo -e "${YELLOW}🏥 Health check for all services:${NC}"
    
    local services=("photo-svc:8001" "user-svc:8003" "transaction-svc:8005" "upload-svc:8002" "notification-svc:8004")
    
    for service in "${services[@]}"; do
        local service_name=$(echo $service | cut -d: -f1)
        local port=$(echo $service | cut -d: -f2)
        
        echo -e "${BLUE}Checking $service_name...${NC}"
        
        # Check if container is running
        if docker ps --filter "name=$service_name" --filter "status=running" | grep -q $service_name; then
            # Check health endpoint
            if curl -s http://localhost:$port/health > /dev/null 2>&1; then
                echo -e "${GREEN}✅ $service_name is healthy${NC}"
            else
                echo -e "${RED}❌ $service_name health check failed${NC}"
            fi
        else
            echo -e "${RED}❌ $service_name is not running${NC}"
        fi
    done
}

# Function to show network information
show_network() {
    echo -e "${YELLOW}🌐 Network information:${NC}"
    
    echo -e "${BLUE}Networks:${NC}"
    docker network ls
    
    echo -e "${BLUE}Port mappings:${NC}"
    docker ps --format "table {{.Names}}\t{{.Ports}}"
    
    echo -e "${BLUE}Network connections:${NC}"
    netstat -tulpn | grep -E "(8001|8002|8003|8004|8005|5432|6379|4222|8500|9000|9001)" | head -20
}

# Function to show storage information
show_storage() {
    echo -e "${YELLOW}💾 Storage information:${NC}"
    
    echo -e "${BLUE}Docker volumes:${NC}"
    docker volume ls
    
    echo -e "${BLUE}Volume usage:${NC}"
    docker system df -v
}

# Function to show environment variables
show_env() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Environment variables for $service:${NC}"
    docker exec $service env | sort
}

# Function to execute command in container
exec_command() {
    local service=$1
    local command=${2:-bash}
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Executing '$command' in $service:${NC}"
    docker exec -it $service $command
}

# Function to show service logs
show_service_logs() {
    local service=$1
    local lines=${2:-100}
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}📋 Logs for $service (last $lines lines):${NC}"
    docker logs --tail=$lines $service
}

# Function to restart service
restart_service() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Available services: ${SERVICES[*]}${NC}"
        return
    fi
    
    echo -e "${YELLOW}🔄 Restarting $service...${NC}"
    docker restart $service
    
    # Wait for service to be ready
    echo -e "${YELLOW}⏳ Waiting for $service to be ready...${NC}"
    sleep 10
    
    # Check if service is running
    if docker ps --filter "name=$service" --filter "status=running" | grep -q $service; then
        echo -e "${GREEN}✅ $service restarted successfully${NC}"
    else
        echo -e "${RED}❌ $service failed to restart${NC}"
    fi
}

# Function to show help
show_help() {
    echo -e "${BLUE}Usage: $0 [command] [options]${NC}"
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  status                    - Show overall status"
    echo -e "  services [service]        - Show service details"
    echo -e "  resources                 - Show resource usage"
    echo -e "  logs [service] [lines]     - Show logs (default: 50 lines)"
    echo -e "  follow [service]          - Follow logs in real-time"
    echo -e "  events [service]          - Show recent events"
    echo -e "  health                    - Check health endpoints"
    echo -e "  network                   - Show network information"
    echo -e "  storage                   - Show storage information"
    echo -e "  env [service]             - Show environment variables"
    echo -e "  exec [service] [command]  - Execute command in container"
    echo -e "  restart [service]         - Restart service"
    echo -e "  all                       - Show comprehensive status"
    echo -e "  help                      - Show this help message"
    echo -e "${YELLOW}Examples:${NC}"
    echo -e "  $0 status                - Show overall status"
    echo -e "  $0 logs photo-svc 100    - Show last 100 lines of photo-svc logs"
    echo -e "  $0 follow user-svc        - Follow user-svc logs"
    echo -e "  $0 health                - Check all service health"
    echo -e "  $0 exec photo-svc bash   - Open bash in photo-svc container"
    echo -e "  $0 restart user-svc      - Restart user-svc"
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
    show_events
    echo
    check_health
}

# Main function
main() {
    check_docker
    check_docker_compose
    
    case "${1:-status}" in
        "status")
            show_status
            ;;
        "services")
            show_service_details $2
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
        "env")
            show_env $2
            ;;
        "exec")
            exec_command $2 $3
            ;;
        "restart")
            restart_service $2
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
