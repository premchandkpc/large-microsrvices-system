FROM rust:1.76-slim-bookworm AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock* ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release 2>/dev/null || true

COPY . .
RUN cargo build --release --bin search-server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/search-server /search-server

EXPOSE 8086
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8086/health || exit 1

ENTRYPOINT ["/search-server"]
