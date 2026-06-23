#!/bin/bash
set -euo pipefail

echo "=== Platform Health Check ==="

check_service() {
  local name=$1
  local url=$2
  local status=$(curl -s -o /dev/null -w "%{http_code}" --max-time 3 "$url" 2>/dev/null || echo "000")
  if [ "$status" = "200" ] || [ "$status" = "401" ]; then
    echo "  ✅ $name ($url) - $status"
  else
    echo "  ❌ $name ($url) - $status"
  fi
}

echo ""
echo "HTTP Services:"
check_service "API Gateway" "http://localhost:8081/health"
check_service "Auth Service" "http://localhost:8082/health"
check_service "User Service" "http://localhost:8083/health"
check_service "Document Ingestion" "http://localhost:8084/health"
check_service "Document Processing" "http://localhost:8085/health"
check_service "Search Service" "http://localhost:8086/health"
check_service "Notification Service" "http://localhost:8087/health"
check_service "Orchestration Service" "http://localhost:8088/health"
check_service "Analytics Service" "http://localhost:8089/health"
check_service "Frontend" "http://localhost:3001"

echo ""
echo "Infrastructure:"
check_service "Prometheus" "http://localhost:9090/-/ready"
check_service "Grafana" "http://localhost:3000/api/health"
check_service "Elasticsearch" "http://localhost:9200/_cluster/health"
check_service "MinIO" "http://localhost:9000/minio/health/live"
check_service "Qdrant" "http://localhost:6333/healthz"

echo ""
echo "Message Broker:"
echo "  Kafka topics:"
docker exec platform-kafka kafka-topics --list --bootstrap-server localhost:9092 2>/dev/null || echo "  Kafka not available"

echo ""
echo "Databases:"
echo -n "  PostgreSQL: "
docker exec platform-postgres pg_isready -U platform 2>/dev/null && echo "ready" || echo "not available"
echo -n "  Redis: "
docker exec platform-redis redis-cli ping 2>/dev/null || echo "not available"
