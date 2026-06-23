# Document Intelligence Platform

A **polyglot microservices** system for document ingestion, AI/LLM processing, vector search, workflow orchestration, and real-time notifications. Built for production-scale document intelligence.

```
Languages:   Go  |  Java (Spring Boot)  |  Python (FastAPI)  |  Rust (Actix)
Protocols:   REST  |  gRPC  |  Kafka Events  |  WebSocket
Infra:       Docker  |  K8s  |  AWS (EKS/RDS/MSK/ElastiCache)
Data:        PostgreSQL  |  Redis  |  Elasticsearch  |  Qdrant  |  S3/MinIO
AI/LLM:      OpenAI  |  Ollama  |  LangChain  |  Vector Embeddings  |  RAG
```

---

## Architecture Overview

```
[Frontend (React/TS)]
    |
[API Gateway (Go)] ———→ [Auth Service (Java)] ———→ PostgreSQL
    |                      [User Service (Java)] ———→ PostgreSQL
    |                      [Document Ingestion (Go)] —→ S3/MinIO → Kafka
    |                      [Document Processing (Python)] → Qdrant / ES / LLM
    |                      [Search Service (Rust)] ———→ Qdrant + ES
    |                      [Notification Service (Go)] → WebSocket / Email / Push
    |                      [Orchestration Service (Python)] → Kafka / DB
    |                      [Analytics Service (Java)] ———→ Kafka → PostgreSQL
    |
    └── Kafka (Event Bus) ←——→ All Services
```

## Service Map

| Service | Language | Stack | Port | Purpose |
|---------|----------|-------|------|---------|
| `api-gateway` | Go | Gin, gRPC, Kafka, Redis, OTel | 8081 | Routing, auth, rate-limit, WS |
| `auth-service` | Java | Spring Boot, JPA, JWT, Flyway | 8082/50051 | AuthN/AuthZ, JWT, RBAC |
| `user-service` | Java | Spring Boot, JPA, Flyway | 8083 | User profiles, preferences |
| `document-ingestion` | Go | Gin, S3 SDK, Kafka | 8084 | File upload, S3 storage, events |
| `document-processing` | Python | FastAPI, LangChain, Qdrant, ES | 8085 | Chunking, embedding, LLM, indexing |
| `search-service` | Rust | Actix, Qdrant, ES | 8086 | Vector + fulltext + hybrid search |
| `notification-service` | Go | Gin, Kafka, WebSocket, SMTP | 8087 | Multi-channel notifications |
| `orchestration-service` | Python | FastAPI, SQLAlchemy, Kafka | 8088 | Workflow engine, step orchestration |
| `analytics-service` | Java | Spring Boot, Kafka, JPA | 8089 | Event tracking, reporting |
| `frontend` | React/TS | React Router, Axios, Chart.js | 3000 | Web UI |

## Data Stores

| Store | Type | Services | Purpose |
|-------|------|----------|---------|
| PostgreSQL 16 | Relational (RDS) | Auth, User, Analytics | Users, roles, tokens, analytics |
| Redis 7 | Cache (ElastiCache) | Gateway, Ingestion, Notifications | Sessions, rate-limit, cache |
| Elasticsearch 8 | Full-text search | Processing, Search | Full-text document search |
| Qdrant | Vector DB | Processing, Search | Embedding storage, vector search |
| MinIO / S3 | Object storage | Ingestion, Processing | Document file storage |
| Kafka | Event stream (MSK) | All services | Async events, messaging |

## Messaging (Kafka Topics)

| Topic | Partitions | Producer | Consumer | Purpose |
|-------|-----------|----------|----------|---------|
| `document-processing` | 3 | Document Ingestion | Document Processing | Trigger AI processing |
| `document-processed` | 3 | Document Processing | Search, Orchestration | Processing results |
| `notifications` | 3 | All services | Notification Service | Send notifications |
| `analytics-events` | 3 | All services | Analytics Service | Usage tracking |

## Quick Start

```bash
# Prerequisites: Docker, Docker Compose

# 1. Clone and set up
git clone https://github.com/premchandkpc/large-microsrvices-system.git
cd large-microsrvices-system

# 2. Start everything
make setup

# 3. Seed sample data
make seed

# 4. Check health
make health
```

## Local URLs

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3001 |
| API Gateway | http://localhost:8081 |
| Auth Service | http://localhost:8082 |
| Kafka UI | http://localhost:8088 |
| MinIO Console | http://localhost:9001 |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| MailHog | http://localhost:8025 |

Default login: `admin@platform.local` / `admin123`

## Development

```bash
# Run specific service locally
cd services/api-gateway
go run ./cmd/main.go

# Run Java service
cd services/auth-service
mvn spring-boot:run

# Run Python service
cd services/document-processing
uvicorn app.main:app --reload

# Run Rust service
cd services/search-service
cargo run
```

## Production Deployment

See `terraform/` for AWS infrastructure (EKS, RDS, MSK, ElastiCache, S3) and `k8s/` for Kubernetes manifests with HPA, probes, and secrets management.

```bash
# Deploy with Terraform
cd terraform
terraform init
terraform plan
terraform apply

# Deploy with kubectl
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmaps/
kubectl apply -f k8s/secrets/
kubectl apply -f k8s/services/
```

## CI/CD

GitHub Actions pipeline (`.github/workflows/ci.yml`):
1. **Lint** — vet, compile-check across all 4 languages
2. **Test** — unit + integration with testcontainers
3. **Build** — docker images → ECR
4. **Security** — Trivy vulnerability scan
5. **Deploy** — Helm upgrade to EKS

## Observability

- **Metrics**: Prometheus + Grafana dashboards
- **Tracing**: OpenTelemetry → Jaeger / Tempo
- **Logging**: JSON structured logs → Loki
- **Health**: /health, /ready endpoints with probes
- **Alerting**: Prometheus AlertManager rules

## Key Architecture Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| API Gateway | Go (Gin) | High throughput, low latency routing |
| Auth | Java (Spring Security) | Mature RBAC, JWT ecosystem |
| Document processing | Python (FastAPI) | Best LLM/AI library ecosystem |
| Search | Rust (Actix) | Maximum performance for vector search |
| Stateful services | Java (Spring Boot + JPA) | Strong consistency needs |
| Event bus | Kafka | Durable, replayable, partitioned events |
| Vector DB | Qdrant | Purpose-built, high-performance ANN |
| Search | ES + Qdrant hybrid | Best of keyword + semantic search |
| Config | Env vars + K8s secrets | 12-factor app principles |
| Deployment | K8s (EKS) + Helm | Portability, scalability |

## License

MIT
