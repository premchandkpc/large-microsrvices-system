from pydantic_settings import BaseSettings
from typing import List


class Settings(BaseSettings):
    service_name: str = "document-processing"
    port: int = 8085
    environment: str = "development"

    kafka_brokers: str = "localhost:9092"
    kafka_input_topic: str = "document-processing"
    kafka_output_topic: str = "document-processed"
    kafka_analytics_topic: str = "document-analytics"
    kafka_group_id: str = "document-processing"

    redis_addr: str = "localhost:6379"

    qdrant_host: str = "localhost"
    qdrant_port: int = 6333
    qdrant_collection: str = "documents"

    elasticsearch_url: str = "http://localhost:9200"
    elasticsearch_index: str = "documents"

    openai_api_key: str = ""
    openai_model: str = "gpt-4-turbo-preview"
    embedding_model: str = "text-embedding-3-small"

    ollama_base_url: str = "http://localhost:11434"
    ollama_model: str = "llama2"
    use_ollama: bool = False

    storage_type: str = "s3"
    s3_endpoint: str = "http://localhost:9000"
    s3_access_key: str = ""
    s3_secret_key: str = ""
    s3_region: str = "us-east-1"
    s3_bucket: str = "documents"

    max_chunk_size: int = 1000
    chunk_overlap: int = 200
    max_retries: int = 3

    class Config:
        env_prefix = "DP_"
        case_sensitive = False


settings = Settings()
