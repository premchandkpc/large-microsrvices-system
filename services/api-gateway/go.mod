module github.com/premchandkpc/large-microsrvices-system/services/api-gateway

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/websocket v1.5.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.17.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/sony/gobreaker v0.5.0
	github.com/spf13/viper v1.18.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.46.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)
