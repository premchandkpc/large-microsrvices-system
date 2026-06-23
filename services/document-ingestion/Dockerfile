FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /document-ingestion ./cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /document-ingestion /document-ingestion

EXPOSE 8084
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1

ENTRYPOINT ["/document-ingestion"]
