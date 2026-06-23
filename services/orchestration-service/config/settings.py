from pydantic_settings import BaseSettings
from typing import List


class Settings(BaseSettings):
    service_name: str = "orchestration-service"
    port: int = 8088
    environment: str = "development"

    kafka_brokers: str = "localhost:9092"
    kafka_group_id: str = "orchestration-service"

    redis_addr: str = "localhost:6379"

    postgres_url: str = "postgresql+asyncpg://platform:platform_secret_2024@localhost:5432/platform_db"

    api_gateway_url: str = "http://api-gateway:8081"
    document_processing_url: str = "http://document-processing:8085"
    notification_service_url: str = "http://notification-service:8087"

    max_concurrent_workflows: int = 50
    workflow_timeout_seconds: int = 3600

    class Config:
        env_prefix = "ORCH_"
        case_sensitive = False


settings = Settings()
