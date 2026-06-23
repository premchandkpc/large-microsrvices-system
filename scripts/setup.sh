#!/bin/bash
set -euo pipefail

echo "=== Setting up Document Intelligence Platform ==="

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Docker required"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose required"; exit 1; }

# Create .env if missing
if [ ! -f .env ]; then
  cat > .env << 'EOF'
JWT_SECRET=local-dev-jwt-secret-change-in-production-32chars
OPENAI_API_KEY=
SMTP_HOST=mailhog
SMTP_PORT=1025
EOF
  echo "Created .env file"
fi

# Create required directories
mkdir -p infra/postgres/init

# Create MinIO bucket
cat > infra/minio/init.sh << 'SCRIPT'
#!/bin/sh
sleep 5
mc alias set local http://minio:9000 platformadmin platform_secret_2024
mc mb local/documents --ignore-existing
mc policy set download local/documents
SCRIPT

echo "=== Starting infrastructure services ==="
docker-compose up -d postgres redis elasticsearch qdrant minio kafka zookeeper

echo "=== Waiting for services to be ready ==="
echo "Waiting for PostgreSQL..."
until docker exec platform-postgres pg_isready -U platform 2>/dev/null; do
  sleep 2
done
echo "PostgreSQL ready"

echo "Waiting for Kafka..."
until docker exec platform-kafka kafka-topics --list --bootstrap-server localhost:9092 2>/dev/null; do
  sleep 3
done
echo "Kafka ready"

# Create Kafka topics
echo "Creating Kafka topics..."
docker exec platform-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic document-processing \
  --partitions 3 \
  --replication-factor 1 || true

docker exec platform-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic document-processed \
  --partitions 3 \
  --replication-factor 1 || true

docker exec platform-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic notifications \
  --partitions 3 \
  --replication-factor 1 || true

docker exec platform-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic analytics-events \
  --partitions 3 \
  --replication-factor 1 || true

echo "=== Starting application services ==="
docker-compose up -d auth-service user-service document-ingestion api-gateway

echo "=== Building and starting AI services ==="
docker-compose up -d document-processing search-service notification-service orchestration-service

echo "=== Starting frontend ==="
docker-compose up -d frontend

echo "=== Starting monitoring ==="
docker-compose up -d prometheus grafana

echo ""
echo "=== Platform is running ==="
echo "  Frontend:      http://localhost:3001"
echo "  API Gateway:   http://localhost:8081"
echo "  Auth Service:  http://localhost:8082"
echo "  Kafka UI:      http://localhost:8088"
echo "  MinIO Console: http://localhost:9001"
echo "  Grafana:       http://localhost:3000"
echo "  Prometheus:    http://localhost:9090"
echo "  MailHog:       http://localhost:8025"
echo ""
echo "Default admin: admin@platform.local / admin123"
echo ""
echo "To stop all services: docker-compose down"
echo "To view logs: docker-compose logs -f [service-name]"
