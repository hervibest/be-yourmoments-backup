#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Setting up Be Your Moments Development Environment${NC}"

# Function to check if Docker is running
check_docker() {
    if ! docker info &> /dev/null; then
        echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker is running${NC}"
}

# Function to check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose is not available. Please install Docker Compose.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker Compose is available${NC}"
}

# Function to create development directory
create_dev_directory() {
    echo -e "${YELLOW}📁 Creating development directory...${NC}"
    
    # Create development directory
    mkdir -p /tmp/be-yourmoments-dev
    
    # Copy necessary files
    cp docker-compose.yaml /tmp/be-yourmoments-dev/
    cp docker-compose-development.yaml /tmp/be-yourmoments-dev/
    cp -r init/ /tmp/be-yourmoments-dev/
    cp -r scripts/ /tmp/be-yourmoments-dev/
    
    echo -e "${GREEN}✅ Development directory created${NC}"
}

# Function to create environment files
create_env_files() {
    echo -e "${YELLOW}📝 Creating environment files...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Create .env file for infrastructure
    cat > .env << 'EOF'
# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres_password

# Redis Configuration
REDIS_PASSWORD=redis_password

# MinIO Configuration
MINIO_ROOT_USER=minio_access_key
MINIO_ROOT_PASSWORD=minio_secret_key
EOF

    # Create service-specific environment files
    cat > .env.photo << 'EOF'
PORT=8001
DB_HOST=postgres
DB_PORT=5432
DB_USER=photo_user
DB_PASSWORD=photo_password
DB_NAME=photo_db
REDIS_HOST=redis
REDIS_PASSWORD=redis_password
NATS_URL=nats://nats:4222
CONSUL_HOST=consul:8500
EOF

    cat > .env.user << 'EOF'
PORT=8003
DB_HOST=postgres
DB_PORT=5432
DB_USER=user_user
DB_PASSWORD=user_password
DB_NAME=user_db
REDIS_HOST=redis
REDIS_PASSWORD=redis_password
NATS_URL=nats://nats:4222
CONSUL_HOST=consul:8500
EOF

    cat > .env.transaction << 'EOF'
PORT=8005
DB_HOST=postgres
DB_PORT=5432
DB_USER=transaction_user
DB_PASSWORD=transaction_password
DB_NAME=transaction_db
REDIS_HOST=redis
REDIS_PASSWORD=redis_password
NATS_URL=nats://nats:4222
CONSUL_HOST=consul:8500
EOF

    cat > .env.upload << 'EOF'
PORT=8002
DB_HOST=postgres
DB_PORT=5432
DB_USER=upload_user
DB_PASSWORD=upload_password
DB_NAME=upload_db
MINIO_ENDPOINT=http://minio:9000
MINIO_ACCESS_KEY=minio_access_key
MINIO_SECRET_KEY=minio_secret_key
NATS_URL=nats://nats:4222
CONSUL_HOST=consul:8500
EOF

    cat > .env.notification << 'EOF'
PORT=8004
DB_HOST=postgres
DB_PORT=5432
DB_USER=notification_user
DB_PASSWORD=notification_password
DB_NAME=notification_db
REDIS_HOST=redis
REDIS_PASSWORD=redis_password
NATS_URL=nats://nats:4222
CONSUL_HOST=consul:8500
EOF

    echo -e "${GREEN}✅ Environment files created${NC}"
}

# Function to create development docker-compose file
create_dev_compose() {
    echo -e "${YELLOW}🐳 Creating development docker-compose file...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    cat > docker-compose-dev.yaml << 'EOF'
version: '3.8'

services:
  photo-svc:
    image: ghcr.io/hervipro/be-yourmoments-backup-photo-svc:dev
    env_file:
      - .env.photo
    ports:
      - "8001:8001"
    networks:
      - backend
    depends_on:
      - postgres
      - consul
      - nats
    restart: unless-stopped

  transaction-svc:
    image: ghcr.io/hervipro/be-yourmoments-backup-transaction-svc:dev
    env_file:
      - .env.transaction
    ports:
      - "8005:8005"
    networks:
      - backend
    depends_on:
      - redis
      - postgres
      - consul
      - nats
    restart: unless-stopped

  upload-svc:
    image: ghcr.io/hervipro/be-yourmoments-backup-upload-svc:dev
    env_file:
      - .env.upload
    ports:
      - "8002:8002"
    networks:
      - backend
    depends_on:
      - minio
      - consul
      - nats
    restart: unless-stopped

  user-svc:
    image: ghcr.io/hervipro/be-yourmoments-backup-user-svc:dev
    env_file:
      - .env.user
    volumes:
      - ./user-svc/serviceAccountKey.json:/app/serviceAccountKey.json
    environment:
      GOOGLE_APPLICATION_CREDENTIALS: /app/serviceAccountKey.json
    ports:
      - "8003:8003"
    networks:
      - backend
    depends_on:
      - redis
      - postgres
      - consul
      - nats
    restart: unless-stopped

  notification-svc:
    image: ghcr.io/hervipro/be-yourmoments-backup-notification-svc:dev
    env_file:
      - .env.notification
    volumes:
      - ./notification-svc/serviceAccountKey.json:/app/serviceAccountKey.json
    environment:
      GOOGLE_APPLICATION_CREDENTIALS: /app/serviceAccountKey.json
    ports:
      - "8004:8004"
    networks:
      - backend
    depends_on:
      - redis
      - postgres
      - consul
      - nats
    restart: unless-stopped

  redis:
    image: redis:7.2-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    restart: unless-stopped
    networks:
      - backend
    ports:
      - "6379:6379"

  postgres:
    image: postgres:17-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init:/docker-entrypoint-initdb.d
    networks:
      - backend
    ports:
      - "5432:5432"

  minio:
    image: minio/minio:RELEASE.2025-04-22T22-12-26Z
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio-data:/data
    networks:
      - backend
    restart: unless-stopped

  nats:
    image: nats:2.11.2-alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    networks:
      - backend
    command: ["-js", "-m", "8222"]
    restart: unless-stopped

  consul:
    image: hashicorp/consul:1.20.6
    ports:
      - "8500:8500"
    command: "agent -dev -client=0.0.0.0"
    networks:
      - backend
    restart: unless-stopped

networks:
  backend:
    name: backend
    driver: bridge

volumes:
  pgdata:
  minio-data:
EOF

    echo -e "${GREEN}✅ Development docker-compose file created${NC}"
}

# Function to create service account keys
create_service_accounts() {
    echo -e "${YELLOW}🔑 Creating service account keys...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Create directories for service account keys
    mkdir -p user-svc notification-svc
    
    # Create dummy service account keys (replace with real ones)
    cat > user-svc/serviceAccountKey.json << 'EOF'
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "your-private-key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_PRIVATE_KEY_HERE\n-----END PRIVATE KEY-----\n",
  "client_email": "your-service-account@your-project.iam.gserviceaccount.com",
  "client_id": "your-client-id",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project.iam.gserviceaccount.com"
}
EOF

    cp user-svc/serviceAccountKey.json notification-svc/serviceAccountKey.json
    
    echo -e "${GREEN}✅ Service account keys created${NC}"
    echo -e "${YELLOW}⚠️  Please replace with your actual service account keys${NC}"
}

# Function to create management scripts
create_management_scripts() {
    echo -e "${YELLOW}📜 Creating management scripts...${NC}"
    
    cd /tmp/be-yourmoments-dev
    
    # Create start script
    cat > start.sh << 'EOF'
#!/bin/bash
echo "🚀 Starting Be Your Moments Development Environment"
docker-compose -f docker-compose-dev.yaml up -d
echo "✅ Services started"
echo "📊 Service Status:"
docker-compose -f docker-compose-dev.yaml ps
EOF

    # Create stop script
    cat > stop.sh << 'EOF'
#!/bin/bash
echo "🛑 Stopping Be Your Moments Development Environment"
docker-compose -f docker-compose-dev.yaml down
echo "✅ Services stopped"
EOF

    # Create restart script
    cat > restart.sh << 'EOF'
#!/bin/bash
echo "🔄 Restarting Be Your Moments Development Environment"
docker-compose -f docker-compose-dev.yaml down
docker-compose -f docker-compose-dev.yaml up -d
echo "✅ Services restarted"
EOF

    # Create logs script
    cat > logs.sh << 'EOF'
#!/bin/bash
echo "📋 Showing logs for all services"
docker-compose -f docker-compose-dev.yaml logs -f
EOF

    # Create status script
    cat > status.sh << 'EOF'
#!/bin/bash
echo "📊 Service Status:"
docker-compose -f docker-compose-dev.yaml ps
echo ""
echo "💾 Resource Usage:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
EOF

    # Make scripts executable
    chmod +x *.sh
    
    echo -e "${GREEN}✅ Management scripts created${NC}"
}

# Function to show next steps
show_next_steps() {
    echo -e "${GREEN}🎉 Development environment setup completed!${NC}"
    echo -e "${BLUE}📋 Next steps:${NC}"
    echo -e "1. Update service account keys in user-svc/ and notification-svc/"
    echo -e "2. Start services: cd /tmp/be-yourmoments-dev && ./start.sh"
    echo -e "3. Check status: ./status.sh"
    echo -e "4. View logs: ./logs.sh"
    echo -e "5. Stop services: ./stop.sh"
    echo ""
    echo -e "${YELLOW}🌐 Service URLs:${NC}"
    echo -e "Photo Service: http://localhost:8001"
    echo -e "User Service: http://localhost:8003"
    echo -e "Transaction Service: http://localhost:8005"
    echo -e "Upload Service: http://localhost:8002"
    echo -e "Notification Service: http://localhost:8004"
    echo -e "MinIO Console: http://localhost:9001"
    echo -e "Consul UI: http://localhost:8500"
    echo ""
    echo -e "${YELLOW}🔧 Management Commands:${NC}"
    echo -e "cd /tmp/be-yourmoments-dev"
    echo -e "./start.sh    # Start all services"
    echo -e "./stop.sh     # Stop all services"
    echo -e "./restart.sh  # Restart all services"
    echo -e "./status.sh   # Show service status"
    echo -e "./logs.sh     # Show service logs"
}

# Main function
main() {
    check_docker
    check_docker_compose
    create_dev_directory
    create_env_files
    create_dev_compose
    create_service_accounts
    create_management_scripts
    show_next_steps
}

# Run main function
main "$@"
