#!/bin/bash
set -euo pipefail

echo "=== Seeding Platform Data ==="

# Create admin user via Auth Service
echo "Creating admin user..."
curl -s -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@platform.local",
    "password": "admin123",
    "name": "System Admin"
  }' || echo "Admin may already exist"

# Get admin token
echo "Getting admin token..."
TOKEN=$(curl -s -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@platform.local","password":"admin123"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
  echo "Failed to get auth token"
  exit 1
fi

echo "Admin token obtained"

# Create sample workflows via Orchestration Service
echo "Creating sample workflow..."
curl -s -X POST http://localhost:8088/api/v1/workflows \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Sample Document Processing",
    "steps": [
      {"type": "document_upload", "params": {}},
      {"type": "document_process", "params": {"pipeline": "full"}},
      {"type": "index_document", "params": {}},
      {"type": "notify_user", "params": {
        "title": "Document Processed",
        "message": "Your document has been processed successfully"
      }}
    ],
    "user_id": "00000000-0000-0000-0000-000000000001"
  }'

echo ""
echo "=== Seed complete ==="
