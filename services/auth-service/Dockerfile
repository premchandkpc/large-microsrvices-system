FROM eclipse-temurin:17-jdk-alpine AS builder
WORKDIR /app
COPY pom.xml ./
COPY src ./src
RUN apk add --no-cache maven && mvn package -DskipTests -q

FROM eclipse-temurin:17-jre-alpine
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/target/*.jar /app/service.jar

EXPOSE 50051 8082
HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

ENTRYPOINT ["java", "-jar", "/app/service.jar"]
