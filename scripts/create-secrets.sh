#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

NAMESPACE="be-yourmoments"

echo -e "${BLUE}🔐 Creating secrets for Be Your Moments microservices${NC}"

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is not installed or not in PATH${NC}"
        exit 1
    fi
}

# Function to create secret interactively
create_secret_interactive() {
    local secret_name=$1
    local service_name=$2
    
    echo -e "${YELLOW}📝 Creating secrets for $service_name${NC}"
    
    # Database configuration
    read -p "Database host (default: postgres-service): " db_host
    db_host=${db_host:-postgres-service}
    
    read -p "Database port (default: 5432): " db_port
    db_port=${db_port:-5432}
    
    read -p "Database user: " db_user
    read -s -p "Database password: " db_password
    echo
    
    read -p "Database name: " db_name
    
    # Redis configuration
    read -p "Redis host (default: redis-service): " redis_host
    redis_host=${redis_host:-redis-service}
    
    read -s -p "Redis password: " redis_password
    echo
    
    # NATS configuration
    read -p "NATS URL (default: nats://nats-service:4222): " nats_url
    nats_url=${nats_url:-nats://nats-service:4222}
    
    # Consul configuration
    read -p "Consul host (default: consul-service:8500): " consul_host
    consul_host=${consul_host:-consul-service:8500}
    
    # Service-specific configurations
    if [ "$service_name" = "upload-svc" ]; then
        read -p "MinIO endpoint (default: http://minio-service:9000): " minio_endpoint
        minio_endpoint=${minio_endpoint:-http://minio-service:9000}
        
        read -p "MinIO access key: " minio_access_key
        read -s -p "MinIO secret key: " minio_secret_key
        echo
    fi
    
    # Create the secret
    kubectl create secret generic $secret_name \
        --from-literal=db-host="$db_host" \
        --from-literal=db-port="$db_port" \
        --from-literal=db-user="$db_user" \
        --from-literal=db-password="$db_password" \
        --from-literal=db-name="$db_name" \
        --from-literal=redis-host="$redis_host" \
        --from-literal=redis-password="$redis_password" \
        --from-literal=nats-url="$nats_url" \
        --from-literal=consul-host="$consul_host" \
        -n $NAMESPACE
    
    # Add service-specific secrets
    if [ "$service_name" = "upload-svc" ]; then
        kubectl patch secret $secret_name -n $NAMESPACE --type='json' -p='[
            {"op": "add", "path": "/data/minio-endpoint", "value": "'$(echo -n "$minio_endpoint" | base64)'"},
            {"op": "add", "path": "/data/minio-access-key", "value": "'$(echo -n "$minio_access_key" | base64)'"},
            {"op": "add", "path": "/data/minio-secret-key", "value": "'$(echo -n "$minio_secret_key" | base64)'"}
        ]'
    fi
    
    echo -e "${GREEN}✅ Secret $secret_name created successfully${NC}"
}

# Function to create service account secrets
create_service_account_secrets() {
    echo -e "${YELLOW}🔑 Creating service account secrets${NC}"
    
    # User service account
    if [ -f "user-svc/serviceAccountKey.json" ]; then
        kubectl create secret generic user-svc-service-account \
            --from-file=serviceAccountKey.json=user-svc/serviceAccountKey.json \
            -n $NAMESPACE
        echo -e "${GREEN}✅ User service account secret created${NC}"
    else
        echo -e "${YELLOW}⚠️  user-svc/serviceAccountKey.json not found, skipping user service account${NC}"
    fi
    
    # Notification service account
    if [ -f "notification-svc/serviceAccountKey.json" ]; then
        kubectl create secret generic notification-svc-service-account \
            --from-file=serviceAccountKey.json=notification-svc/serviceAccountKey.json \
            -n $NAMESPACE
        echo -e "${GREEN}✅ Notification service account secret created${NC}"
    else
        echo -e "${YELLOW}⚠️  notification-svc/serviceAccountKey.json not found, skipping notification service account${NC}"
    fi
}

# Function to create all secrets
create_all_secrets() {
    local services=("photo-svc" "user-svc" "transaction-svc" "upload-svc" "notification-svc")
    
    for service in "${services[@]}"; do
        local secret_name="${service}-secrets"
        
        # Check if secret already exists
        if kubectl get secret $secret_name -n $NAMESPACE &> /dev/null; then
            echo -e "${YELLOW}⚠️  Secret $secret_name already exists. Skipping...${NC}"
            continue
        fi
        
        create_secret_interactive $secret_name $service
        echo
    done
    
    create_service_account_secrets
}

# Function to create secrets from environment variables
create_secrets_from_env() {
    echo -e "${YELLOW}🌍 Creating secrets from environment variables${NC}"
    
    # Check required environment variables
    local required_vars=("DB_HOST" "DB_PORT" "DB_USER" "DB_PASSWORD" "DB_NAME" "REDIS_HOST" "REDIS_PASSWORD" "NATS_URL" "CONSUL_HOST")
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            echo -e "${RED}❌ Environment variable $var is not set${NC}"
            exit 1
        fi
    done
    
    # Create secrets for all services
    local services=("photo-svc" "user-svc" "transaction-svc" "upload-svc" "notification-svc")
    
    for service in "${services[@]}"; do
        local secret_name="${service}-secrets"
        
        kubectl create secret generic $secret_name \
            --from-literal=db-host="$DB_HOST" \
            --from-literal=db-port="$DB_PORT" \
            --from-literal=db-user="$DB_USER" \
            --from-literal=db-password="$DB_PASSWORD" \
            --from-literal=db-name="$DB_NAME" \
            --from-literal=redis-host="$REDIS_HOST" \
            --from-literal=redis-password="$REDIS_PASSWORD" \
            --from-literal=nats-url="$NATS_URL" \
            --from-literal=consul-host="$CONSUL_HOST" \
            -n $NAMESPACE
        
        # Add MinIO secrets for upload-svc
        if [ "$service" = "upload-svc" ] && [ -n "$MINIO_ENDPOINT" ] && [ -n "$MINIO_ACCESS_KEY" ] && [ -n "$MINIO_SECRET_KEY" ]; then
            kubectl patch secret $secret_name -n $NAMESPACE --type='json' -p='[
                {"op": "add", "path": "/data/minio-endpoint", "value": "'$(echo -n "$MINIO_ENDPOINT" | base64)'"},
                {"op": "add", "path": "/data/minio-access-key", "value": "'$(echo -n "$MINIO_ACCESS_KEY" | base64)'"},
                {"op": "add", "path": "/data/minio-secret-key", "value": "'$(echo -n "$MINIO_SECRET_KEY" | base64)'"}
            ]'
        fi
        
        echo -e "${GREEN}✅ Secret $secret_name created${NC}"
    done
    
    create_service_account_secrets
}

# Function to show existing secrets
show_secrets() {
    echo -e "${YELLOW}📋 Existing secrets in namespace $NAMESPACE:${NC}"
    kubectl get secrets -n $NAMESPACE
}

# Function to delete secrets
delete_secrets() {
    echo -e "${YELLOW}🗑️  Deleting all secrets in namespace $NAMESPACE${NC}"
    kubectl delete secrets --all -n $NAMESPACE
    echo -e "${GREEN}✅ All secrets deleted${NC}"
}

# Function to show help
show_help() {
    echo -e "${BLUE}Usage: $0 [command]${NC}"
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  interactive  - Create secrets interactively (default)"
    echo -e "  env          - Create secrets from environment variables"
    echo -e "  show         - Show existing secrets"
    echo -e "  delete        - Delete all secrets"
    echo -e "  help         - Show this help message"
    echo -e "${YELLOW}Environment Variables (for 'env' command):${NC}"
    echo -e "  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME"
    echo -e "  REDIS_HOST, REDIS_PASSWORD"
    echo -e "  NATS_URL, CONSUL_HOST"
    echo -e "  MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY (for upload-svc)"
}

# Main function
main() {
    check_kubectl
    
    case "${1:-interactive}" in
        "interactive")
            create_all_secrets
            ;;
        "env")
            create_secrets_from_env
            ;;
        "show")
            show_secrets
            ;;
        "delete")
            delete_secrets
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
