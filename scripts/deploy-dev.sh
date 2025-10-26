#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TAG=${1:-dev}
COMPOSE_FILE="docker-compose.yaml"
BACKUP_DIR="/tmp/be-yourmoments-backup-$(date +%Y%m%d-%H%M%S)"

echo -e "${BLUE}🚀 Deploying Be Your Moments Development Environment${NC}"
echo -e "${YELLOW}Tag: $TAG${NC}"

# Function to check if Docker is running
check_docker() {
    if ! docker info &> /dev/null; then
        echo -e "${RED}❌ Docker is not running${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker is running${NC}"
}

# Function to check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not available${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker Compose is available${NC}"
}

# Function to backup current deployment
backup_current() {
    echo -e "${YELLOW}📦 Creating backup of current deployment...${NC}"
    
    if [ -d "/tmp/be-yourmoments-dev" ]; then
        cp -r /tmp/be-yourmoments-dev $BACKUP_DIR
        echo -e "${GREEN}✅ Backup created at $BACKUP_DIR${NC}"
    else
        echo -e "${YELLOW}⚠️  No existing deployment to backup${NC}"
    fi
}

# Function to stop existing services
stop_services() {
    echo -e "${YELLOW}🛑 Stopping existing services...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Stop services gracefully
    if command -v docker-compose &> /dev/null; then
        docker-compose down --remove-orphans || true
    else
        docker compose down --remove-orphans || true
    fi
    
    echo -e "${GREEN}✅ Services stopped${NC}"
}

# Function to pull latest images
pull_images() {
    echo -e "${YELLOW}📥 Pulling latest images...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Update image tags in docker-compose.yaml
    sed -i "s/:latest/:$TAG/g" $COMPOSE_FILE
    
    # Pull images
    if command -v docker-compose &> /dev/null; then
        docker-compose pull
    else
        docker compose pull
    fi
    
    echo -e "${GREEN}✅ Images pulled successfully${NC}"
}

# Function to start services
start_services() {
    echo -e "${YELLOW}🚀 Starting services...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Start services
    if command -v docker-compose &> /dev/null; then
        docker-compose up -d
    else
        docker compose up -d
    fi
    
    echo -e "${GREEN}✅ Services started${NC}"
}

# Function to wait for services to be ready
wait_for_services() {
    echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        echo -e "${YELLOW}Attempt $attempt/$max_attempts...${NC}"
        
        # Check if all services are running
        if command -v docker-compose &> /dev/null; then
            local running_services=$(docker-compose ps --services --filter "status=running" | wc -l)
            local total_services=$(docker-compose ps --services | wc -l)
        else
            local running_services=$(docker compose ps --services --filter "status=running" | wc -l)
            local total_services=$(docker compose ps --services | wc -l)
        fi
        
        if [ "$running_services" -eq "$total_services" ] && [ "$total_services" -gt 0 ]; then
            echo -e "${GREEN}✅ All services are running${NC}"
            return 0
        fi
        
        sleep 10
        attempt=$((attempt + 1))
    done
    
    echo -e "${RED}❌ Services failed to start within expected time${NC}"
    return 1
}

# Function to check service health
check_health() {
    echo -e "${YELLOW}🏥 Checking service health...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Show service status
    if command -v docker-compose &> /dev/null; then
        docker-compose ps
    else
        docker compose ps
    fi
    
    # Check individual service health
    local services=("photo-svc:8001" "user-svc:8003" "transaction-svc:8005" "upload-svc:8002" "notification-svc:8004")
    
    for service in "${services[@]}"; do
        local service_name=$(echo $service | cut -d: -f1)
        local port=$(echo $service | cut -d: -f2)
        
        echo -e "${YELLOW}Checking $service_name...${NC}"
        
        # Wait for service to be ready
        local max_attempts=10
        local attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if curl -s http://localhost:$port/health > /dev/null 2>&1; then
                echo -e "${GREEN}✅ $service_name is healthy${NC}"
                break
            fi
            
            if [ $attempt -eq $max_attempts ]; then
                echo -e "${RED}❌ $service_name health check failed${NC}"
            fi
            
            sleep 5
            attempt=$((attempt + 1))
        done
    done
}

# Function to show logs
show_logs() {
    echo -e "${YELLOW}📋 Recent logs:${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    if command -v docker-compose &> /dev/null; then
        docker-compose logs --tail=20
    else
        docker compose logs --tail=20
    fi
}

# Function to cleanup old images
cleanup_images() {
    echo -e "${YELLOW}🧹 Cleaning up old images...${NC}"
    
    # Remove unused images
    docker image prune -f
    
    # Remove old versions of our images
    docker images | grep "be-yourmoments" | grep -v "$TAG" | awk '{print $3}' | xargs -r docker rmi -f
    
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Function to show deployment status
show_status() {
    echo -e "${BLUE}📊 Deployment Status:${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    echo -e "${YELLOW}Services:${NC}"
    if command -v docker-compose &> /dev/null; then
        docker-compose ps
    else
        docker compose ps
    fi
    
    echo -e "${YELLOW}Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    
    echo -e "${YELLOW}Ports:${NC}"
    echo "Photo Service: http://localhost:8001"
    echo "User Service: http://localhost:8003"
    echo "Transaction Service: http://localhost:8005"
    echo "Upload Service: http://localhost:8002"
    echo "Notification Service: http://localhost:8004"
    echo "MinIO Console: http://localhost:9001"
    echo "Consul UI: http://localhost:8500"
}

# Function to rollback
rollback() {
    echo -e "${YELLOW}🔄 Rolling back to previous deployment...${NC}"
    
    if [ -d "$BACKUP_DIR" ]; then
        # Stop current services
        stop_services
        
        # Restore backup
        rm -rf /tmp/be-yourmoments-dev
        cp -r $BACKUP_DIR /tmp/be-yourmoments-dev
        
        # Start services
        start_services
        
        echo -e "${GREEN}✅ Rollback completed${NC}"
    else
        echo -e "${RED}❌ No backup found for rollback${NC}"
        exit 1
    fi
}

# Function to show help
show_help() {
    echo -e "${BLUE}Usage: $0 [command] [options]${NC}"
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  deploy [tag]     - Deploy with specific tag (default: dev)"
    echo -e "  status           - Show deployment status"
    echo -e "  logs             - Show service logs"
    echo -e "  health           - Check service health"
    echo -e "  rollback         - Rollback to previous deployment"
    echo -e "  cleanup          - Cleanup old images"
    echo -e "  stop             - Stop all services"
    echo -e "  restart          - Restart all services"
    echo -e "  help             - Show this help message"
}

# Function to stop services
stop_all() {
    echo -e "${YELLOW}🛑 Stopping all services...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    if command -v docker-compose &> /dev/null; then
        docker-compose down
    else
        docker compose down
    fi
    
    echo -e "${GREEN}✅ All services stopped${NC}"
}

# Function to restart services
restart_services() {
    echo -e "${YELLOW}🔄 Restarting services...${NC}"
    
    stop_all
    start_services
    wait_for_services
    
    echo -e "${GREEN}✅ Services restarted${NC}"
}

# Main deployment function
deploy() {
    echo -e "${GREEN}🎯 Starting deployment process${NC}"
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Backup current deployment
    backup_current
    
    # Stop existing services
    stop_services
    
    # Pull latest images
    pull_images
    
    # Start services
    start_services
    
    # Wait for services to be ready
    wait_for_services
    
    # Check health
    check_health
    
    # Show status
    show_status
    
    echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
}

# Main function
main() {
    case "${1:-deploy}" in
        "deploy")
            deploy
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs
            ;;
        "health")
            check_health
            ;;
        "rollback")
            rollback
            ;;
        "cleanup")
            cleanup_images
            ;;
        "stop")
            stop_all
            ;;
        "restart")
            restart_services
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
