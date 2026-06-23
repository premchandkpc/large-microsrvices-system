FROM eclipse-temurin:17-jdk-alpine AS builder
WORKDIR /app
COPY pom.xml ./
COPY src ./src
RUN apk add --no-cache maven && mvn package -DskipTests -q

FROM eclipse-temurin:17-jre-alpine
COPY --from=builder /app/target/*.jar /app/service.jar
EXPOSE 8083
ENTRYPOINT ["java", "-jar", "/app/service.jar"]
